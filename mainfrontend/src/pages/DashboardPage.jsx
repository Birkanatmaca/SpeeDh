import React from 'react';
import { useNavigate } from 'react-router-dom';
import { BsFillMicFill, BsUpload, BsClockHistory } from "react-icons/bs";
import './DashboardPage.scss';

const DashboardPage = () => {
    const navigate = useNavigate();
    const userName = "Kullanıcı"; // Bu bilgi daha sonra backend'den alınacak

    return (
        <div className="dashboard-page">
            {/* Arka Plan İkonları */}
            <div className="background-icons">
                <div className="background-icon-wrapper icon1"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon2"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon3"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon4"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon5"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon6"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon7"><BsFillMicFill /></div>
            </div>

            {/* İçerik Kutusu */}
            <div className="dashboard-box">
                <header className="dashboard-header">
                    <h1>SpeeDh'e hoş geldin, {userName}!</h1>
                    <button className="history-button" onClick={() => navigate('/history')}>
                        <BsClockHistory /> Geçmiş Transkriptler
                    </button>
                </header>

                <main className="dashboard-content">
                    {/* Ses Kaydedici Alanı */}
                    <div className="recorder-section">
                        <h2>Ses Kaydet</h2>
                        <div className="recorder-visualizer">
                            <BsFillMicFill className="mic-icon" />
                            <div className="record-time">00:00</div>
                        </div>
                        <button className="record-button">Kaydı Başlat</button>
                    </div>

                    {/* Dosya Yükleme Alanı */}
                    <div className="upload-section">
                        <h2>Veya Bir Dosya Yükle</h2>
                        <div className="upload-area">
                            <BsUpload className="upload-icon" />
                            <p>Dosyanızı buraya sürükleyin veya seçmek için tıklayın.</p>
                            <input type="file" className="file-input" />
                        </div>
                        <button className="transcribe-button">Metne Dönüştür</button>
                    </div>
                </main>
            </div>
        </div>
    );
};

export default DashboardPage;