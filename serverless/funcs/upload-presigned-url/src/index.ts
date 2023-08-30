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

// ////////////////////////////////////////////////////////////////////////////
// exports
// ////////////////////////////////////////////////////////////////////////////
export const handler: Handler = async ( {
  queryStringParameters,
}: APIGatewayEvent ): Promise<APIGatewayProxyResult> => {
  const key = nanoid( 24 );

  // eslint-disable-next-line dot-notation
  const rawContentType = queryStringParameters?.['contentType'];
  if( !rawContentType ) {
    return {
        statusCode: 400,
        headers: {
          'content-type': 'application/json',
        },
        body: JSON.stringify( {
          error: 'No content type submitted'
        } ),
      };
  }

  const contentType = decodeURIComponent(rawContentType);

  const s3Params = {
    Bucket: S3_UPLOAD_BUCKET,
    Key: key,
    Expires: URL_EXPIRATION_SECONDS || 300,
    ContentType: contentType,
  };

  const uploadURL = await s3.getSignedUrlPromise( 'putObject', s3Params );

  return {
    statusCode: 201,
    headers: {
      'content-type': 'application/json',
      'Access-Control-Allow-Headers': 'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token',
      'Access-Control-Allow-Methods': 'GET,POST,OPTIONS',
      'Access-Control-Allow-Origin':  '*',
    },
    body: JSON.stringify( {
      uploadURL,
      key,
    } ),
  };
};
