// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import prettyBytes from 'pretty-bytes';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { submitFiles } from './upload';
import { showError, showSuccess, showWarning } from './alert';

import currentUser from '../stores/current-user';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
const MEDIA_TYPES = [
  'application',
  'audio',
  'font',
  'image',
  'text',
  'video',
] as const;

type TMediaType = typeof MEDIA_TYPES[number];

type TAllowedExtensions = string[];
type TFileSubtypeMap = Record<string, TAllowedExtensions>;
type TFileTypeMap = Record<TMediaType, TFileSubtypeMap>;

// ////////////////////////////////////////////////////////////////////////////
// Constants
// ////////////////////////////////////////////////////////////////////////////
const MAX_FILE_SIZE = 1000 * 1000 * 1000; // 1GB

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
let fileToUpload: Nullable<File> = null;

// NB: Some files may be represented by multiple MIME types and/or extensions
const FILE_VALIDATION_MAP: TFileTypeMap = {
  application: {
    'epub+zip': ['epub'],
    msword: ['doc'],
    pdf: ['pdf'],
    postscript: [
      'ai', 'eps', 'ps',
    ],
    psd: ['psd'],
    rle: ['rle'],
    rtf: ['rtf'],
    scitex: ['sct'],
    'vnd.openxmlformats-officedocument.presentationml.presentation': ['pptx'],
    'vnd.openxmlformats-officedocument.spreadsheetml.sheet': ['xlsx'],
    'vnd.openxmlformats-officedocument.wordprocessingml.document': ['docx'],
    'vnd.ms-excel': ['xls'],
    'vnd.ms-powerpoint': ['ppt'],
    'vnd.rar': ['rar'],
    'vnd.Quark.QuarkXPress': ['qxd', 'qxp'],
    'x-indesign': ['indd'],
    'x-rle': ['rle'],
    'x-shockwave-flash': ['swf'],
    'x-subrip': ['srt'],
    zip: ['zip'],
  },
  audio: {
    aiff: ['aif', 'aiff'],
    mp3: ['mp3'],
    mp4: ['m4a'],
    mpeg: ['mp3'],
    wav: ['wav'],
    'x-aiff': ['aif', 'aiff'],
    'x-m4a': ['m4a'],
    'x-ms-wma': ['wma'],
  },
  font: {
    otf: ['otf'],
    ttf: ['ttf'],
  },
  image: {
    bmp: ['bmp'],
    emf: ['emf'],
    gif: ['gif'],
    jpeg: ['jpeg', 'jpg'],
    png: ['png'],
    psd: ['psd'],
    rle: ['rle'],
    'svg+xml': ['svg'],
    tiff: ['tif', 'tiff'],
    'vnd.adobe.photoshop': ['psb', 'psd'],
    'vnd.fpx': ['fpx'],
    'vnd.zbrush.pcx': ['pcx'],
    webp: ['webp'],
    'x-adobe-dng': ['dng'],
    'x-canon-cr2': ['cr2'],
    'x-canon-crw': ['crw'],
    'x-dcx': ['dcx'],
    'x-emf': ['emf'],
    'x-fuji-raf': ['raf'],
    'x-icon': ['ico'],
    'x-kodak-dcr': ['dcr'],
    'x-minolta-mrw': ['mrw'],
    'x-nikon-nef': ['nef', 'nrw'],
    'x-olympus-orf': ['orf'],
    'x-pcx': ['pcx'],
    'x-pict': ['pct', 'pic'],
    'x-pentax-pef': ['pef'],
    'x-photo-cd': ['pcd'],
    'x-sigma-x3f': ['x3f'],
    'x-sony-arw': ['arw'],
    'x-sun-raster': ['ras'],
    'x-tga': ['tga'],
    'x-wmf': ['wmf'],
    'x-wpg': ['wpg'],
  },
  text: {
    csv: ['csv'],
    html: [
      'htm', 'html', 'shtml', 'xhtml',
    ], // Is anything other than ".html" allowed by Aprimo?
    plain: ['txt'],
    vtt: ['vtt'],
  },
  video: {
    avi: ['avi'],
    dv: ['dv'],
    mp4: ['m4v', 'mp4'],
    mpeg: ['mpg', 'mpeg'],
    mxf: ['mxf'],
    quicktime: ['mov', 'qt'], // Is anything other than ".mov" allowed by Aprimo?
    webm: ['webm'],
    'x-dv': ['dv'],
    'x-flv': ['flv'],
    'x-m4v': ['m4v'],
    'x-ms-wmv': ['wmv'],
    'x-msvideo': ['avi'],
  },
};

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
const fileExtension = ( fileName: string ) => {
  const segments = fileName.split( '.' );

  return segments[segments.length - 1];
};

const typeInfo = ( fileType: string ) => {
  const segments = fileType.split( '/' );

  return { mediaType: segments[0], subtype: segments[1] };
};

const validateFile = ( { type, size, name }: File ) => {
  const ext = fileExtension( name );
  const { mediaType, subtype } = typeInfo( type );

  if ( !ext || ext.length < 2 || ext.length > 4 ) {
    showError( 'Uploaded files must have a valid file extension' );

    return false;
  }

  const lookupType = FILE_VALIDATION_MAP[mediaType as TMediaType];

  if ( !lookupType ) {
    showError( 'Invalid media type. Most audio, video, picture, text, and Office documents are supported.' );

    return false;
  }

  const allowedExtensions = lookupType[subtype];

  if ( !allowedExtensions ) {
    showError( 'Invalid file type. Try converting to a more common file format.' );

    return false;
  }

  const extensionMatched = allowedExtensions.includes( ext );

  if ( !extensionMatched ) {
    showError( 'Invalid file extension. Make sure your file extension matches the file type, such as ".docx" for a Word file.' );

    return false;
  }

  if ( size > MAX_FILE_SIZE ) {
    showError( `Max file size is ${prettyBytes( MAX_FILE_SIZE )}, but the uploaded file is ${prettyBytes( size )}` );

    return false;
  }

  return true;
};

const setUpload = ( file: File ) => {
  if ( !validateFile( file ) ) {
    return;
  }

  const list = document.getElementById( 'file-list' ) as HTMLElement;

  list.innerHTML = `${file.name} (${prettyBytes( file.size )})`;

  fileToUpload = file;
};

const handleFile = ( files?: Nullable<FileList> ) => {
  if ( files && files.length > 1 ) {
    showError( 'Only single-file uploads are currently supported' );
  } else if ( files ) {
    setUpload( files[0] );
  }
};

const validateSubmission = ( descriptionElem: HTMLInputElement ) => {
  const description = descriptionElem?.value;
  const file = fileToUpload;
  const { email, team } = currentUser.get();

  let error = false;

  // Error Check
  if ( !descriptionElem ) {
    showError( 'Internal error' );
    error = true;
  } else if ( !file ) {
    showWarning( 'No file has been selected for upload', 'Invalid Submission' );
    error = true;
  } else if ( !description ) {
    showWarning( 'No file description provided', 'Invalid Submission' );
    error = true;
  } else if ( !email ) {
    showError( 'No current user email' );
    error = true;
  } else if ( !team ) {
    showError( 'No current user team' );
    error = true;
  }

  return { description, file, error, email, team };
};

const switchSubmitDisplay = () => {
  const btn = document.getElementById( 'upload-files-btn' ) as HTMLInputElement;
  const btnHidden = ( btn.style.display === 'none' );

  btn.style.display = ( btnHidden ? 'block' : 'none' );

  const loader = document.getElementById( 'loader' ) as HTMLInputElement;
  const loaderHidden = ( loader.style.display === 'none' );

  loader.style.display = ( loaderHidden ? 'block' : 'none' );
};

// ////////////////////////////////////////////////////////////////////////////
// Exports
// ////////////////////////////////////////////////////////////////////////////

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

/**
 * Prepares the drag and dropped files for upload.
 *
 * @param e The drop event.
 */
export const dropHandler = ( e: DragEvent ) => {
  haltEvent( e );

  const files = e?.dataTransfer?.files;

  handleFile( files );
};

export const chooseHandler = ( e: Event ) => {
  const { files } = ( e.target as HTMLInputElement );

  handleFile( files );
};

export const submitHandler = async () => {
  switchSubmitDisplay();

  // Prepare and validate
  const descriptionElem = document.getElementById( 'description-text' ) as HTMLInputElement;
  const fileElem = document.getElementById( 'file-list' ) as HTMLInputElement;

  const { file, description, email, team, error } = validateSubmission( descriptionElem );

  if ( error ) {
    switchSubmitDisplay();

    return;
  }

  // Send data
  const response = await submitFiles( file as File, { description, email, team } );

  if ( response !== 'ok' ) {
    showError( 'Could not upload file' );
  } else {
    showSuccess( 'File upload has been initiated' );
  }

  // Cleanup
  descriptionElem.value = '';
  fileElem.innerHTML = '';
  fileToUpload = null;

  switchSubmitDisplay();
};
