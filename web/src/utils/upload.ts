// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { buildQuery } from './api';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
type TUploadStatus = 'ok' | 'urlFailed' | 'uploadFailed' | 'metadataFailed';

interface IMetadata {
  description: string;
  email?: string;
  team?: string;
}

interface IFileUploadMeta extends IMetadata {
  key: string;
  fileType: string;
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const UPLOAD_ENDPOINT = import.meta.env.UPLOAD_ENDPOINT || 'upload';

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const getPresignedUrl = async ( filename: string, contentType: string ) => {
  const queryFilename = `fileName=${encodeURIComponent( filename )}`;
  const queryType = `contentType=${encodeURIComponent( contentType )}`;

  const response = await buildQuery( `${UPLOAD_ENDPOINT}?${queryFilename}&${queryType}`, null, 'GET' );
  const { uploadURL, key } = await response.json();

  return { uploadURL, key };
};

const uploadFile = async ( fqUrl: string, file: File ) => {
  const response = await fetch( fqUrl, {
    body: file,
    method: 'PUT',
  } );

  return response.status === 200;
};

const submitFileMetadata = async ( body: IFileUploadMeta ) => {
  const response = await buildQuery( UPLOAD_ENDPOINT, body, 'POST' );

  return response.status === 200;
};

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////
export const submitFiles = async ( file: File, meta: IMetadata ): Promise<TUploadStatus> => {
  const { uploadURL, key } = await getPresignedUrl( file.name, file.type );

  if ( !uploadURL ) { return 'urlFailed'; }

  const fileSuccess = await uploadFile( uploadURL, file );

  if ( !fileSuccess ) { return 'uploadFailed'; }

  const metaSuccess = await submitFileMetadata( { ...meta, key, fileType: file.type } );

  if ( !metaSuccess ) { return 'metadataFailed'; }

  return 'ok';
};
