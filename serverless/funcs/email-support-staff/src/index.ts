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

interface ITeam {
  teamName: string,
  teamId: string,
}

type WithTeam<T> = T & ITeam

interface ISupportStaffEventBody {
  contentCommonsUser: IUser,
  externalTeamLead: WithTeam<IUser>,
  supportStaff: IUser[],
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
const ses = new SESClient( { region: AWS_SES_REGION || '' } );

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
function formatEmailBody( { contentCommonsUser, externalTeamLead, supportStaff }: ISupportStaffEventBody ) {
  return `\
<p>Hello ${contentCommonsUser.givenName} ${contentCommonsUser.familyName},</p>

<p>${externalTeamLead.givenName} ${externalTeamLead.familyName} from team ${externalTeamLead.teamName}
 has requested ${supportStaff.length} support staff be authorized to their team.
  Please review the list below and approve ordeny access at your convenience.</p>
<ul>
  ${
  supportStaff.map( user => `<li>${user.givenName} ${user.familyName} (${user.email})</li>` ).join( '' )
}
</ul>
<p>Thank you,<br>The Content Commons Team</p>
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

    const recvData: ISupportStaffEventBody = JSON.parse( body );
    const email = formatEmail( recvData );

    const { MessageId: emailMessageId } = await ses.send( email );

    return { eventMessageId, emailMessageId };
  } );

  const results = await Promise.allSettled( promises );

  results.forEach( result => {
    if ( result.status === 'fulfilled' ) {
      const { value: { emailMessageId, eventMessageId } } = result;

      logger.info( `Send email with ID ${emailMessageId} for event ${eventMessageId}` );
    } else {
      logger.error( result.reason );
    }
  } );
};
