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
const API_ENDPOINT = import.meta.env.PUBLIC_SERVERLESS_URL;
const UPLOAD_ENDPOINT = import.meta.env.UPLOAD_ENDPOINT || 'upload';

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const getPresignedUrl = async ( contentType: string ) => {
  const response = await fetch( `${API_ENDPOINT}/${UPLOAD_ENDPOINT}?contentType=${encodeURIComponent( contentType )}`, {
    method: 'GET',
  } );
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
  const response = await fetch( `${API_ENDPOINT}/${UPLOAD_ENDPOINT}`, {
    body: JSON.stringify( body ),
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
  } );

  return response.status === 200;
};

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////
export const submitFiles = async ( file: File, meta: IMetadata ): Promise<TUploadStatus> => {
  const { uploadURL, key } = await getPresignedUrl( file.type );

  if ( !uploadURL ) { return 'urlFailed'; }

  const fileSuccess = await uploadFile( uploadURL, file );

  if ( !fileSuccess ) { return 'uploadFailed'; }

  const metaSuccess = await submitFileMetadata( { ...meta, key, fileType: file.type } );

  if ( !metaSuccess ) { return 'metadataFailed'; }

  return 'ok';
};
