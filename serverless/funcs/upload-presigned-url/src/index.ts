// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { Handler, APIGatewayEvent, APIGatewayProxyResult } from 'aws-lambda';
import { config, S3 } from 'aws-sdk';

import { nanoid } from 'nanoid';

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const {
  AWS_REGION,
  S3_UPLOAD_BUCKET,
  URL_EXPIRATION_SECONDS,
} = process.env;

config.update( { region: AWS_REGION || 'us-east-1' } );

const s3 = new S3();

const CORS_HEADERS = {

};

// ////////////////////////////////////////////////////////////////////////////
// exports
// ////////////////////////////////////////////////////////////////////////////
export const handler: Handler = async ( { queryStringParameters }: APIGatewayEvent ): Promise<APIGatewayProxyResult> => {
  const key = nanoid(24);
  const contentType = queryStringParameters.contentType;

  const s3Params = {
    Bucket: S3_UPLOAD_BUCKET,
    Key: key,
    Expires: URL_EXPIRATION_SECONDS,
    ContentType: contentType,
  };

  const uploadURL = await s3.getSignedUrlPromise( 'putObject', s3Params );

  return {
    statusCode: 201,
    headers: CORS_HEADERS,
    body: JSON.stringify( {
      uploadURL,
      key,
    } ),
  }
};
