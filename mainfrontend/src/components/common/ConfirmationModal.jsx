import React from 'react';
import './ConfirmationModal.scss';

const ConfirmationModal = ({ isOpen, onClose, onConfirm, message }) => {
    if (!isOpen) {
        return null;
    }

    return (
        <div className="confirmation-modal-overlay">
            <div className="confirmation-modal-content">
                <p className="confirmation-modal-message">{message}</p>
                <div className="confirmation-modal-actions">
                    <button className="confirm-button yes" onClick={onConfirm}>Evet</button>
                    <button className="confirm-button no" onClick={onClose}>Hayır</button>

                </div>
            </div>
        </div>
    );
};

export default ConfirmationModal;