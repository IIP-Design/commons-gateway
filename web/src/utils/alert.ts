import Swal from 'sweetalert2';

export const showError = ( text: string ) => Swal.fire( {
  icon: 'error',
  title: 'Input Error',
  text,
} );

export const showSuccess = ( text: string ) => Swal.fire( {
  icon: 'success',
  title: 'Success!',
  text,
} );

export const showConfirm = ( text: string ) => Swal.fire( {
  icon: 'question',
  title: 'Please confirm your selection',
  text,
  showCancelButton: true,
  showConfirmButton: true,
} );
