import prettyBytes from 'pretty-bytes';
import Swal from 'sweetalert2';
import { submitFiles } from './api';

const MAX_FILE_SIZE = 100 * 1000 * 1000;

const formData = new FormData();

const showError = ( text: string ) => {
  return Swal.fire( {
    icon: 'error',
    title: 'Input Error',
    text,
  } );
}

const showSuccess = ( text: string ) => {
  return Swal.fire( {
    icon: 'success',
    title: 'Success!',
    text,
  } );
}

/**
   * Prevents the default browser behavior (i.e. opening the file) when
   * a file is dropped into the browser.
   *
   * @param e The dragenter/dragover event.
   */
export const haltEvent = ( e: Event ) => {
  e.stopPropagation();
  e.preventDefault();
};

const validateFile = ( { type, size }: File ) => {
  if( !type.match( /^(image|video)\/.+/ ) ) {
    showError('Only pictures and/or videos may be uploaded');
    return false;
  } else if( size > MAX_FILE_SIZE ) {
    showError(`Max file size is ${prettyBytes(MAX_FILE_SIZE)}`);
    return false;
  } else {
    return true;
  }
}

const addToUploadList = ( file: File ) => {
  if( !validateFile(file) ) {
    return;
  }

  const list = document.getElementById( 'file-list' );
  const listItem = document.createElement( 'li' );

  console.log( file.size );

  listItem.innerHTML = file.name;
  list?.appendChild( listItem );

  formData.append('file', file);
};

const handleFiles = ( files: FileList ) => {
  [...files].forEach( file => addToUploadList( file ) );
};

/**
 * Prepares the drag and dropped files for upload.
 *
 * @param e The drop event.
 */
export const dropHandler = ( e: DragEvent ) => {
  haltEvent(e);

  console.log( e );

  const files = e?.dataTransfer?.files;

  if ( files ) {
    handleFiles( files );
  }
};

export const submitHandler = async () => {
  const descriptionElem = document.getElementById( 'description-text' ) as HTMLInputElement;
  const description = descriptionElem?.value;

  const listElem = document.getElementById( 'file-list' );
  const listEntries = listElem?.childElementCount;

  if( !descriptionElem || !listElem ) {
    showError( 'Internal error' );
    return;
  } else if( !listEntries ) {
    showError('No files have been selected for upload');
    return;
  } else if( !description ) {
    showError('No file description provided');
    return;
  }

  formData.append( 'description', description );

  const response = await submitFiles( 'upload', formData );
  if( !response.ok ) {
    showError( `Error on the ${response.status >= 500 ? 'server' : 'client'}` );
    return;
  } else {
    showSuccess( `File${listEntries > 1 ? 's' : ''} have been uploaded` );
  }

  descriptionElem.value = '';
  listElem.innerHTML = '';
  formData.delete('file');
  formData.delete('description');
}