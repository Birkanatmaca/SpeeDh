import React from 'react';
import { BsX, BsDownload } from 'react-icons/bs';
import './Modal.scss';

const Modal = ({ isOpen, onClose, transcript, onDownloadText, onDownloadAudio }) => {
    if (!isOpen || !transcript) {
        return null;
    }

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                    <h3 className="modal-title">{transcript.title}</h3>
                    <button className="modal-close-button" onClick={onClose}>
                        <BsX />
                    </button>
                </div>
                <div className="modal-body">
                    <p>{transcript.transcription_text}</p>
                </div>
                <div className="modal-footer">
                    <button className="modal-action-button" onClick={onDownloadText}>
                        <BsDownload /> Transkripti İndir (.txt)
                    </button>
                    <button className="modal-action-button" onClick={onDownloadAudio}>
                        <BsDownload /> Sesi İndir
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Modal;