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

interface ISupportStaffEventBody {
  contentCommonsUser: IUser,
  externalTeamLead: IUser,
  supportStaffUser: IUser,
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

const logger = new Logger( { serviceName: AWS_SERVICE_NAME || 'email-support-staff' } );
const ses = new SESClient( { region: AWS_SES_REGION || 'us-east-1' } );

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
function formatEmailBody( { contentCommonsUser, externalTeamLead, supportStaffUser, url }: ISupportStaffEventBody ) {
  return `\
<p>${contentCommonsUser.givenName} ${contentCommonsUser.familyName},</p>

<p>${externalTeamLead.givenName} ${externalTeamLead.familyName} has submitted a ticket for adding
 ${supportStaffUser.givenName} ${supportStaffUser.familyName} for your approval.
  Please follow <a href="${url}">this link</a> to approve or deny this request.</p>
<p>This email was generated automatically. Please do not reply to this email.</p>
`;
}

function formatEmail( recvData: ISupportStaffEventBody ) {
  return new SendEmailCommand( {
    Destination: {
      ToAddresses: [recvData.contentCommonsUser.email],
    },
    Message: {
      Body: {
        Html: {
          Charset: 'UTF-8',
          Data: formatEmailBody( recvData ),
        },
      },
      Subject: {
        Data: `Content Commons Support Staff Request`,
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

    const recvData: ISupportStaffEventBody = JSON.parse( body );
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
