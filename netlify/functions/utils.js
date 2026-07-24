// isDevEnv ...
const isDevEnv = process.env.DEV_ENV === "true";

// FacebookEnvValues ...
// defines common facebook env values
export const FacebookEnvValues = {
  FileName: "CONTENT_FILENAME",
  FacebookPageId: "FACEBOOK_PAGE_ID",
  FacebookPageUri: "FACEBOOK_URI",
  FacebookPageAccessToken: "FACEBOOK_PAGE_ACCESS_TOKEN",
};

// validateRequest ...
// defines a function that is used to validate a request
export const validateRequest = (apiKey = "") => {
  if (isDevEnv) return true;
  if (!isDevEnv && apiKey === IntegrationApiKey) return true;
  return false;
};

// populateCorsHeaders ...
// defines a function that populates cors headers for each functions
export const populateCorsHeaders = () => {
  return {
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Methods": "POST, OPTIONS",
    "Access-Control-Allow-Headers": "Content-Type",
  };
};

// Constant ...
// defines the constant values
export const Constant = {
  EmptyDataset: "No data found to process",
  FailedHealthCheck: "Service has failed the health check",
  FailedToPost: "Service has failed to perform post",
};

// ApiConstant ...
// defines the constant responses for Api Requests
export const ApiConstant = {
  HttpStatusOk: "Status Ok",
  HttpStatusBadRequest: "Bad request",
  HttpUnauthorized: "Method not authorized",
  HttpStatusSystemFailed: "Internal server error",
};
