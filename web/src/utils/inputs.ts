/**
 * Switch the input fields from passwords to text when in focus.
 * @param el The input element in question.
 */
export const toggleInputType = ( el: FocusEvent ) => {
  const { type, target } = el;
  const { id } = target as HTMLInputElement;

  const inputType = type === 'focus' ? 'text' : 'password';

  if ( id ) {
    const input = document.getElementById( id ) as HTMLInputElement;

    input.type = inputType;
  }
};
