# Load environment variables from .env file
export $(grep -v '^#' ../.env | xargs)
