import {
  ApiConstant,
  Constant,
  fetchRandomImage,
  populateCorsHeaders,
  populateDataFromGoogleSheets,
  updateGoogleSheetWithStatus,
  validateRequest,
} from "./utils";

export const handler = async (event) => {
  // ARPS validation occurs differently
  const { integrationKey } = JSON.parse(event?.body || "{}");
  if (!integrationKey || integrationKey === "") {
    console.debug(ApiConstant.HttpUnauthorized);
    return {
      statusCode: 401,
      headers: populateCorsHeaders(),
      body: JSON.stringify({ error: ApiConstant.HttpUnauthorized }),
    };
  }

  const isValidRequest = validateRequest(integrationKey);
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
    const results = await populateDataFromGoogleSheets();
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
      URI: process.env.FACEBOOK_URI,
      PageID: process.env.FACEBOOK_PAGE_ID,
      PageToken: process.env.FACEBOOK_PAGE_ACCESS_TOKEN,
    };

    // isFacebookPosted === v[3]
    const dataToPost = results?.find((v) => v[3] === "FALSE");

    if (!dataToPost) {
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

    // if no dataToPost exist; error out before crashing fn
    const toPostRowIdx = results?.findIndex((v) => v[3] === "FALSE");
    const sheetRow = toPostRowIdx + 1;

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

    const imagePath = await fetchRandomImage();

    const isComplete = await performPostToFacebook(
      facebook,
      dataToPost[2],
      imagePath,
    );

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

    // update google sheets after the data is posted
    await updateGoogleSheetWithStatus(sheetRow, "D");

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

// performPost ...
// defines a function that is used to post into facebook
const performPostToFacebook = async (fb, message, imagePath) => {
  const requestUrl = `${fb.URI}/${fb.PageID}/photos?origin_graph_explorer=1&transport=cors&access_token=${fb.PageToken}`;

  const form = new FormData();

  form.append("url", imagePath);
  form.append("message", message);

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
