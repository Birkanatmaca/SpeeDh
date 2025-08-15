import React from 'react';
import { BsFillMicFill } from "react-icons/bs";
import './HistoryPage.scss';

const HistoryPage = () => {
    // Örnek veri
    const transcripts = [
        { id: 1, name: 'toplanti_kaydi_1.wav', date: '14.08.2025', duration: '15:32' },
        { id: 2, name: 'gorusme_notlari.mp3', date: '12.08.2025', duration: '05:48' },
        { id: 3, name: 'podcast_bolum_3.m4a', date: '11.08.2025', duration: '45:12' },
        { id: 4, name: 'fikir_firtinasi.mp3', date: '10.08.2025', duration: '22:05' },
    ];

    return (
        <div className="history-page">
            <div className="background-icons">
                <div className="background-icon-wrapper icon1"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon2"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon3"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon4"><BsFillMicFill /></div>
            </div>

            <div className="history-box">
                <header className="history-header">
                    <h1>Geçmiş Transkriptler</h1>
                    <a href="/dashboard" className="back-button">Ana Sayfaya Dön</a>
                </header>
                <main className="history-content">
                    <table className="history-table">
                        <thead>
                            <tr>
                                <th>Dosya Adı</th>
                                <th>Tarih</th>
                                <th>Süre</th>
                                <th>İşlemler</th>
                            </tr>
                        </thead>
                        <tbody>
                            {transcripts.map(item => (
                                <tr key={item.id}>
                                    <td>{item.name}</td>
                                    <td>{item.date}</td>
                                    <td>{item.duration}</td>
                                    <td>
                                        <button className="action-button view">Görüntüle</button>
                                        <button className="action-button delete">Sil</button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </main>
            </div>
        </div>
    );
};

export default HistoryPage;