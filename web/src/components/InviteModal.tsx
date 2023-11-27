// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useState } from 'react';
import type{ FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import Modal from 'react-modal';
import type { IInvite } from '../utils/types';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import tableStyles from '../styles/table.module.scss';
import btnStyle from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface IModalProps {
  readonly invites: IInvite[];
  readonly anchor: string | JSX.Element;
}

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
const InviteEntry = ( { dateInvited, accessEndDate, pending, expired }: IInvite, idx: number ) => (
  <div key={ `${dateInvited}-${idx}` }>
    <hr style={ { marginTop: '10px' } } />
    <div className="field-group">
      <label>
        <span>Invite Date</span>
        <input
          type="date"
          disabled
          value={ dateInvited }
        />
      </label>
      <label>
        <span>Access End Date</span>
        <input
          type="date"
          disabled
          value={ accessEndDate }
        />
      </label>
    </div>
    <table className={ `${tableStyles.table}` }>
      <thead>
        <tr>
          <th>Status</th>
          <th>Result</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>
            <span className={ tableStyles.status }>
              <span className={ !pending ? tableStyles.active : tableStyles.inactive } />
              { pending ? 'Pending' : 'Approved' }
            </span>
          </td>
          <td>
            <span className={ tableStyles.status }>
              <span className={ !expired ? tableStyles.active : tableStyles.inactive } />
              { expired ? 'Expired' : 'Current' }
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
Modal.setAppElement( document.getElementById( 'root' ) as HTMLElement );

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const InviteModal: FC<IModalProps> = ( { invites, anchor }: IModalProps ) => {
  // Modal Setup
  const [modalIsOpen, setModalIsOpen] = useState( false );
  const noHistory = ( invites.length === 1 );

  // Modal Controls
  const openModal = () => setModalIsOpen( true );
  const closeModal = () => setModalIsOpen( false );

  // Styles
  const modalStyle = {
    content: {
      height: 'fit-content',
      maxHeight: '80%',
      width: '400px',
      maxWidth: '80%',
      top: '50%',
      left: '50%',
      right: 'auto',
      bottom: 'auto',
      marginRight: '-50%',
      padding: '2rem',
      transform: 'translate(-50%, -50%)',
      boxShadow: 'rgba(0, 0, 0, 0.35) 0px 5px 15px',
      overflowY: 'scroll',
    },
  };

  return (
    <>
      <button
        className={ `${btnStyle.btn} ${noHistory ? btnStyle['disabled-btn'] : ''}` }
        onClick={ openModal }
        type="button"
        disabled={ noHistory }
      >
        { anchor }
      </button>
      <Modal
        isOpen={ modalIsOpen }
        onRequestClose={ closeModal }
        contentLabel="Invite History"
        style={ modalStyle }
      >
        <h3>Invite History</h3>
        { invites.slice( 0 ).map( ( invite, idx ) => InviteEntry( invite, idx ) ) }
        <button
          className={ btnStyle.btn }
          onClick={ closeModal }
          type="button"
        >
          Close
        </button>
      </Modal>
    </>
  );
};
