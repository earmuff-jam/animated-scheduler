/**
 * instagramScheduler ...
 * defines a function used to create instagram posts
 *
 * parses single data from csv every day until the marker is turned
 * off using instagram services. Does not have scheduling services
 * to use.
 */

import fs from "node:fs";
import path from "node:path";
import csv from "csv-parser";

import {
  ApiConstant,
  Constant,
  FacebookEnvValues,
  populateCorsHeaders,
  validateRequest,
} from "./utils";
import { createObjectCsvWriter } from "csv-writer";

export const handler = async (event) => {
  const isValidRequest = validateRequest(event.headers["x-api-key"]);
  if (!isValidRequest) {
    console.debug(ApiConstant.HttpUnauthorized);
    return {
      statusCode: 401,
      headers: populateCorsHeaders(),
      body: JSON.stringify({ error: ApiConstant.HttpUnauthorized }),
    };
  }

  if (event.httpMethod !== "POST") {
    console.debug(ApiConstant.HttpUnauthorized);
    return {
      statusCode: 405,
      headers: populateCorsHeaders(),
      body: JSON.stringify({ error: ApiConstant.HttpUnauthorized }),
    };
  }

  try {
    const results = await parseContentsOfCsv();
    if (results.length <= 0) {
      console.debug(Constant.EmptyDataset);
      return {
        statusCode: 500,
        headers: populateCorsHeaders(),
        body: JSON.stringify({
          error: ApiConstant.HttpStatusBadRequest,
          errorDetails: Constant.EmptyDataset,
        }),
      };
    }

    const instagram = {
      URI: process.env[FacebookEnvValues.FacebookPageUri],
      PageID: process.env[FacebookEnvValues.FacebookPageId],
      PageToken: process.env[FacebookEnvValues.FacebookPageAccessToken],
    };

    // retrieves the first dataset that is not posted to instagram
    const dataToPost = results.find(
      (element) => element.IsInstagramComplete !== "true",
    );

    const businessID = await performHealthCheck(instagram);
    if (businessID === "") {
      console.debug(Constant.FailedHealthCheck);
      return {
        statusCode: 500,
        headers: populateCorsHeaders(),
        body: JSON.stringify({
          error: ApiConstant.HttpStatusSystemFailed,
          errorDetails: Constant.FailedHealthCheck,
        }),
      };
    }

    const imagePath = await fetchRandomImage();

    const instagramMediaContainer = await createInstagramMediaContainer(
      businessID,
      instagram,
      dataToPost?.Message,
      imagePath,
    );

    if (instagramMediaContainer === "") {
      console.debug(Constant.FailedToPost);
      return {
        statusCode: 500,
        headers: populateCorsHeaders(),
        body: JSON.stringify({
          error: ApiConstant.HttpStatusSystemFailed,
          errorDetails: Constant.FailedToPost,
        }),
      };
    }

    const isCompleted = postInstagramFromMediaContainer(
      instagramMediaContainer?.id,
      instagram,
      businessID,
    );

    if (!isCompleted) {
      console.debug(Constant.FailedToPost);
      return {
        statusCode: 500,
        headers: populateCorsHeaders(),
        body: JSON.stringify({
          error: ApiConstant.HttpStatusSystemFailed,
          errorDetails: Constant.FailedToPost,
        }),
      };
    }

    // if data was posted, update csv to mark posted as true
    await updateCsvInstagramStatus(results, dataToPost);

    return {
      statusCode: 200,
      headers: populateCorsHeaders(),
      body: JSON.stringify({ message: ApiConstant.HttpStatusOk }),
    };
  } catch (error) {
    console.debug(ApiConstant.HttpStatusSystemFailed, error);
    return {
      statusCode: 500,
      headers: populateCorsHeaders(),
      body: JSON.stringify({
        error: ApiConstant.HttpStatusSystemFailed,
        errorDetails: error,
      }),
    };
  }
};

// postInstagramFromMediaContainer ...
// defines a function that uploads data into instagram
const postInstagramFromMediaContainer = async (
  containerID,
  instagram,
  businessID,
) => {
  if (containerID === "") {
    console.debug("Unable to process without container id.");
    return false;
  }

  const publishUrl = `${instagram.URI}/${businessID}/media_publish?access_token=${instagram.PageToken}`;

  const payload = {
    creation_id: containerID,
  };

  const response = await fetch(publishUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const responseBody = await response.text();

  if (!response.ok) {
    console.debug(
      `instagram media creation failed. status: ${response.status} response: ${responseBody}`,
    );
    return false;
  }

  return true;
};

// createInstagramMediaContainer ...
// defines a function that creates media container for instagram
const createInstagramMediaContainer = async (
  businessID,
  instagram,
  message,
  imageUrl,
) => {
  const mediaUrl = `${instagram.URI}/${businessID}/media?access_token=${instagram.PageToken}`;

  const payload = {
    image_url: imageUrl,
    caption: message,
  };

  const response = await fetch(mediaUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const responseBody = await response.text();

  if (!response.ok) {
    console.debug(
      `instagram media creation failed. status: ${response.status} response: ${responseBody}`,
    );
    return "";
  }

  const result = JSON.parse(responseBody);
  return result;
};

// performHealthCheck ...
// defines a function that is used to perform health check
const performHealthCheck = async (instagram) => {
  const url = `${instagram.URI}/${instagram.PageID}?fields=instagram_business_account&access_token=${instagram.PageToken}`;

  const response = await fetch(url);

  if (!response.ok) {
    const body = await response.text();
    console.debug(`unable to perform health check. Details: ${body}`);
    return "";
  }

  const result = await response.json();
  console.debug("Health check completed. Response:", result);

  return result?.instagram_business_account?.id || "";
};

// parseContentsOfCsv ...
// defines a function that is used to parse the contents of the csv
const parseContentsOfCsv = () => {
  const filenameWithPath = process.env[FacebookEnvValues.FileName];

  return new Promise((resolve, reject) => {
    const results = [];

    fs.createReadStream(filenameWithPath)
      .pipe(csv())
      .on("data", (row) => results.push(row))
      .on("end", () => resolve(results))
      .on("error", reject);
  });
};

// fetchRandomImage ...
// defines a function that is uXsed to fetch a random image
const fetchRandomImage = async () => {
  const response = await fetch("https://picsum.photos/1200/1200");

  if (!response.ok) {
    throw new Error("unable to fetch public image url");
  }

  return response.url;
};

// updateCsvInstagramStatus ...
// defines a function that updates the instagram status for the CSV row
const updateCsvInstagramStatus = async (rows, rowToUpdate) => {
  const filenameWithPath = process.env[FacebookEnvValues.FileName];

  rowToUpdate.IsInstagramComplete = "true";

  const csvWriter = createObjectCsvWriter({
    path: filenameWithPath,
    header: Object.keys(rows[0]).map((key) => ({
      id: key,
      title: key,
    })),
  });

  await csvWriter.writeRecords(rows);
};
