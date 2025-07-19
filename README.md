# Lexia

Golang backend for dictionary application - "Lexia" 

# Translation Service Setup Guide

### 1. Google Cloud Translation API Setup

#### Option A: Service Account (Recommended)
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the Google Translation API:
   ```
   Navigation: APIs & Services > Library > Cloud Translation API > Enable
   ```
4. Create a service account:
   ```
   Navigation: IAM & Admin > Service Accounts > Create Service Account
   ```
5. Add the "Cloud Translation API User" role
6. Generate and download a JSON key file
7. Add to your `.env` file:
   ```properties
   GOOGLE_APPLICATION_CREDENTIALS="C:/path/to/your/service-account-key.json"
   ```

#### Option B: Application Default Credentials (Development)
1. Install Google Cloud CLI
2. Run: `gcloud auth application-default login`
3. No additional environment variables needed

### 2. Get project ID and set it as a environment var `GOOGLE_CLOUD_PROJECT_ID`
