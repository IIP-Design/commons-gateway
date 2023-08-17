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

interface IEmailEventBody {
  invitee: IUser,
  inviter: IUser,
  tmpPassword: string,
  url: string,
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const {
  AWS_SERVICE_NAME,
  AWS_SES_REGION,
  SOURCE_EMAIL_ADDRESS,
} = process.env;

const logger = new Logger( { serviceName: AWS_SERVICE_NAME || 'email-creds' } );
const ses = new SESClient( { region: AWS_SES_REGION || 'us-east-1' } );

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
function formatEmailBody( { invitee, inviter, url, tmpPassword }: IEmailEventBody ) {
  return `\
<p>Hello ${invitee.givenName} ${invitee.familyName},</p>

<p>A new account has been created for you in the Content Commons system by ${inviter.givenName} ${inviter.familyName}.
  Your temporary login information is as follows:</p>
<ul>
  <li>URL: ${url}</li>
  <li>Password: ${tmpPassword}</li>
</ul>
<p>Please login at your convenience to complete your workflow.</p>
<p>Thank you,<br>The Content Commons Team</p>
`;
}

function formatEmail( recvData: IEmailEventBody ) {
  return new SendEmailCommand( {
    Destination: {
      ToAddresses: [recvData.invitee.email],
    },
    Message: {
      Body: {
        Html: {
          Charset: 'UTF-8',
          Data: formatEmailBody( recvData ),
        },
      },
      Subject: {
        Data: 'Content Commons Account Created',
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

    const recvData: IEmailEventBody = JSON.parse( body );
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
