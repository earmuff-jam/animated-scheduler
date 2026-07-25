import fs from "fs";
import path from "path";
import { google } from "googleapis";

// IsDevEnv ...
const IsDevEnv = process.env.DEV_ENV === "true";
const IntegrationApiKey = process.env.INTEGRATION_KEY;

// validateRequest ...
// defines a function that is used to validate a request
export const validateRequest = (apiKey = "") => {
  if (IsDevEnv) return true;
  if (!IsDevEnv && apiKey === IntegrationApiKey) return true;
  return false;
};

// initializeServiceAccount
// returns service account based on dev env
const initializeServiceAccount = (isDevEnv = false) => {
  if (isDevEnv) {
    console.debug("Running in developmental instance. ");
    const serviceAccountPath = path.resolve("./dev/account.json");
    const serviceAccount = JSON.parse(
      fs.readFileSync(serviceAccountPath, "utf8"),
    );

    return new google.auth.GoogleAuth({
      credentials: serviceAccount,
      scopes: ["https://www.googleapis.com/auth/spreadsheets"],
    });
  } else {
    console.log(
      console.log({
        project: process.env.GOOGLE_PROJECT_ID,
        email: process.env.GOOGLE_CLIENT_EMAIL,
        hasKey: !!process.env.GOOGLE_PRIVATE_KEY,
      }),
    );
    return new google.auth.GoogleAuth({
      credentials: {
        type: "service_account",
        project_id: process.env.GOOGLE_PROJECT_ID,
        client_email: process.env.GOOGLE_CLIENT_EMAIL,
        private_key: process.env.GOOGLE_PRIVATE_KEY.replace(/\\n/g, "\n"),
      },
      scopes: ["https://www.googleapis.com/auth/spreadsheets"],
    });
  }
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

// populateDataFromGoogleSheets ...
// defines a function that is used to populate data from google sheets
export const populateDataFromGoogleSheets = async () => {
  const sheets = google.sheets({
    version: "v4",
    auth: initializeServiceAccount(IsDevEnv),
  });

  const response = await sheets.spreadsheets.values.get({
    spreadsheetId: process.env.GOOGLE_SHEET_FILENAME,
    range: "Sheet1",
  });

  const data = response.data;
  return data?.values || [];
};

// updateGoogleSheetWithStatus ...
// defines a function that is used to update google sheets
export const updateGoogleSheetWithStatus = async () => {};

// fetchRandomImage ...
// defines a function that is uXsed to fetch a random image
export const fetchRandomImage = async () => {
  const response = await fetch("https://picsum.photos/1200/1200");

  if (!response.ok) {
    throw new Error("unable to fetch public image url");
  }

  return response.url;
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
