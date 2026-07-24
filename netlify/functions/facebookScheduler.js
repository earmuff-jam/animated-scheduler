/**
 * facebookScheduler ...
 * defines a function used to create facebook scheduled posts
 *
 * Posts all data in CSV for upto 30 days period into facebook using facebook
 * scheduling services.
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

    const facebook = {
      URI: process.env[FacebookEnvValues.FacebookPageUri],
      PageID: process.env[FacebookEnvValues.FacebookPageId],
      PageToken: process.env[FacebookEnvValues.FacebookPageAccessToken],
    };

    const isValid = await performHealthCheck(facebook);
    if (!isValid) {
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

    // post all content into facebook scheduler
    await results.forEach((element) => {
      const imagePath = fetchRandomImage();

      const isComplete = performPostToFacebook(facebook, element, imagePath);
      if (!isComplete) {
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
    });

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

// fetchRandomImage ...
// defines a function that is used to fetch a random image
const fetchRandomImage = () => {
  const imagesDir = path.join("content", "images");

  const entries = fs.readdirSync(imagesDir);

  if (entries.length === 0) {
    console.debug("image directory is empty");
    return "placeholder.png";
  }

  const randomIndex = Math.floor(Math.random() * entries.length);
  const imagePath = path.join(imagesDir, entries[randomIndex]);

  return imagePath;
};

// performPost ...
// defines a function that is used to post into facebook
const performPostToFacebook = async (fb, data, imagePath) => {
  const requestUrl = `${fb.URI}/${fb.PageID}/photos?origin_graph_explorer=1&transport=cors&access_token=${fb.PageToken}`;

  const form = new FormData();

  const image = await fs.openAsBlob(imagePath);

  form.append("source", image, path.basename(imagePath));
  form.append("message", data.message);

  form.append("published", "false");
  form.append(
    "scheduled_publish_time",
    Math.floor(new Date(data.date).getTime() / 1000).toString(),
  );

  const response = await fetch(requestUrl, {
    method: "POST",
    body: form,
  });

  if (!response.ok) {
    throw new Error(await response.text());
  }

  return true;
};

// performHealthCheck ...
// defines a function that is used to perform health check
const performHealthCheck = async (fb) => {
  const url = `${fb.URI}/${fb.PageID}/settings?origin_graph_explorer=1&transport=cors&access_token=${fb.PageToken}`;

  const response = await fetch(url);

  if (response.status === 400) {
    const body = await response.text();
    console.debug(`unable to perform health check. Details: ${body}`);
    return false;
  }

  await response.json();
  return true;
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
