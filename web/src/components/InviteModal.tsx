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
const InviteEntry = ( invite: IInvite ) => {
    return (
      <div key={invite.dateInvited}>
        <div className="field-group">
          <label>
            <span>Invite Date</span>
            <input
              type="date"
              disabled
              value={invite.dateInvited}
            />
          </label>
          <label>
            <span>Access End Date</span>
            <input
              type="date"
              disabled
              value={invite.accessEndDate}
            />
          </label>
        </div>
        <table className={`${tableStyles.table}`}>
          <tbody>
            <tr>
              <td>
              <span className={ tableStyles.status }>
                <span className={ !invite.pending ? tableStyles.active : tableStyles.inactive } />
                { invite.pending ? 'Pending' : 'Approved' }
              </span>
              </td>
              <td>
              <span className={ tableStyles.status }>
                <span className={ !invite.expired ? tableStyles.active : tableStyles.inactive } />
                { invite.expired ? 'Expired' : 'Current' }
              </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    );
  }

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

  // Modal Controls
  const openModal = () => setModalIsOpen( true );
  const closeModal = () => setModalIsOpen( false );

  // Styles
  const modalStyle = {
    content: {
      height: 'fit-content',
      width: '400px',
      minWidth: 'fit-content',
      top: '50%',
      left: '50%',
      right: 'auto',
      bottom: 'auto',
      marginRight: '-50%',
      padding: '2rem',
      transform: 'translate(-50%, -50%)',
      boxShadow: 'rgba(0, 0, 0, 0.35) 0px 5px 15px',
    },
  };

  return (
    <>
      <button className={ btnStyle.btn } onClick={ openModal } type="button">{ anchor }</button>
      <Modal
        isOpen={ modalIsOpen }
        onRequestClose={ closeModal }
        contentLabel="Invite History"
        style={ modalStyle }
      >
        {
            invites.slice(0).map( InviteEntry )
        }
      </Modal>
    </>
  );
};
