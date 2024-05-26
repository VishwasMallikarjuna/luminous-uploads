# Luminous Uploads

Luminous Uploads is an application that provides an HTTP REST API for uploading images via an invitation link. It includes endpoints for generating upload links, uploading images, and retrieving images. The application ensures that each image is only stored once by recognizing duplicates and storing metadata in a PostgreSQL database.

## Database Schema

Create the necessary tables using the following script in `src/db/schema.sql`:

```sql
CREATE TABLE IF NOT EXISTS images
(
    id integer NOT NULL DEFAULT nextval('images_id_seq'::regclass),
    image_data bytea NOT NULL,
    image_hash text COLLATE pg_catalog."default",
    CONSTRAINT images_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS image_detail
(
    id integer NOT NULL,
    width integer NOT NULL,
    height integer NOT NULL,
    camera_model character varying(255) COLLATE pg_catalog."default",
    location character varying(255) COLLATE pg_catalog."default",
    format character varying(50) COLLATE pg_catalog."default" NOT NULL,
    upload_timestamp timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT image_detail_pkey PRIMARY KEY (id)
);

```

Endpoints and Responses
Generate Upload Link

    Endpoint: POST /generate-upload-link/:duration
    Description: Accepts a secret token and expiration time, and produces an expirable link to the endpoint that can be used for image upload.
    Authentication: Bearer Token

    duration (string): The duration for which the upload link is valid.
        10m for 10 minutes
        1h for 1 hour
        24h for 24 hours
        7d for 7 days (Note: d is not supported by time.ParseDuration, so you'd need to convert it to hours, e.g., 168h for 7 days)

    Example Response:

json

{
    "url": "/upload-image",
    "expires_at": "2024-06-01T12:00:00Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Upload Image

    Endpoint: POST /upload-image
    Description: Accepts one or many images and returns back some identifier(s). The logic of uploading multiple images until the link is expired is up to you.
    Authentication: Bearer Token (from POST /generate-upload-link/:duration)
    Example Response:

json

{
[1, 2, 3]
}

Get Image

    Endpoint: GET /image/:imageId
    Description: Accepts an image identifier and returns back an image.
    Authentication: Bearer Token
    Example Response: Returns the image in binary format.



## Configuration

### Azure AD Configuration

To configure Azure AD, follow these guides:

- [Register an application with Azure AD](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
- [Create a service principal in the Azure portal](https://learn.microsoft.com/en-us/entra/identity-platform/howto-create-service-principal-portal)

### Local Configuration

Update the configuration file `config.yml` with your settings:

```yaml
dbUser: postgres-User
dbPassword: postgres-password
dbName: postgres-database-name
dbHost: postgres-localhost
dbPort: postgres-port
clientID: Azure-clientID
tenantID: Azure-tenantID
secretKey: your-secret-key
```

# Running the Application

## Local Setup

To run the application locally, use the following command:

```sh
CONFIG_PATH=path/to/your/config.yml go run main.go
```

## Docker Setup

To pull the Docker image:

```sh
docker pull vishwasmallikarjuna/urlrepo:v0.1
```
To run the Docker image:
```sh
docker run -it -p 1323:1323 -v /path/of/your/local/config.yml:/app/config.yml -e CONFIG_PATH=/app/config.yml vishwasmallikarjuna/urlrepo:v0.1 /bin/sh

```
