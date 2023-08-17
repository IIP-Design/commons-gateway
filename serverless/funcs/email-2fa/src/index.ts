// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import { Handler, SQSEvent } from 'aws-lambda';
import { Logger } from '@aws-lambda-powertools/logger';
import { SESClient, SendEmailCommand } from '@aws-sdk/client-ses';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUser {
  email: string,
  givenName: string,
  familyName: string,
}

interface I2FAEventBody {
  user: IUser,
  verificationCode: string,
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const {
  AWS_SERVICE_NAME,
  AWS_SES_REGION,
  SOURCE_EMAIL_ADDRESS,
} = process.env;

const logger = new Logger( { serviceName: AWS_SERVICE_NAME || 'email-2fa' } );
const ses = new SESClient( { region: AWS_SES_REGION || 'us-east-1' } );

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
function formatEmailBody( { user: { givenName, familyName }, verificationCode }: I2FAEventBody ) {
  return `\
<p>${givenName} ${familyName},</p>

<p>Please use this verification code to complete your sign in:</p>
<p>${verificationCode}</p>
<p>If you did not make this request, please disregard this email. </p>
`;
}

function formatEmail( recvData: I2FAEventBody ) {
  return new SendEmailCommand( {
    Destination: {
      ToAddresses: [recvData.user.email],
    },
    Message: {
      Body: {
        Html: {
          Charset: 'UTF-8',
          Data: formatEmailBody( recvData ),
        },
      },
      Subject: {
        Data: 'Verification Code for Content Commons Login',
      },
    },
    Source: SOURCE_EMAIL_ADDRESS,
  } );
}

// ////////////////////////////////////////////////////////////////////////////
// exports
// ////////////////////////////////////////////////////////////////////////////
export const handler: Handler = async ( { Records: records }: SQSEvent ) => {
  const promises = records.map( async ( { messageId: eventMessageId, body } ) => {
    logger.debug( `Processing event with message ID ${eventMessageId}` );

    const recvData: I2FAEventBody = JSON.parse( body );
    const email = formatEmail( recvData );

    const { MessageId: emailMessageId } = await ses.send( email );

    return { eventMessageId, emailMessageId };
  } );

  const results = await Promise.allSettled( promises );

  results.forEach( result => {
    if ( result.status === 'fulfilled' ) {
      const { value: { emailMessageId, eventMessageId } } = result;

      logger.info( `Sent email with ID ${emailMessageId} for event ${eventMessageId}` );
    } else {
      logger.error( result.reason );
    }
  } );
};
