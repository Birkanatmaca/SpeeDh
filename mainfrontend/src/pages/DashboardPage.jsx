import React, { useState, useRef } from 'react';
import { BsFillMicFill, BsUpload, BsClockHistory, BsGearFill, BsArrowLeft, BsRecordCircle, BsStopCircle } from "react-icons/bs";
import Input from '../components/common/Input';
import Button from '../components/common/Button';
import { transcribeAudio } from '../api/authService'; // API fonksiyonunu import et
import './DashboardPage.scss';

const DashboardPage = () => {
    const [activeView, setActiveView] = useState('main');
    const [isRecording, setIsRecording] = useState(false);
    const [audioBlob, setAudioBlob] = useState(null);
    const [uploadedFile, setUploadedFile] = useState(null);
    const [transcription, setTranscription] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const mediaRecorder = useRef(null);
    const audioChunks = useRef([]);
    const userName = "Kullanıcı";

    // Örnek geçmiş transkript verisi
    const transcripts = [
        { id: 1, name: 'toplanti_kaydi_1.wav', date: '14.08.2025', duration: '15:32' },
        { id: 2, name: 'gorusme_notlari.mp3', date: '12.08.2025', duration: '05:48' },
        { id: 3, name: 'podcast_bolum_3.m4a', date: '11.08.2025', duration: '45:12' },
        { id: 4, name: 'fikir_firtinasi.mp3', date: '10.08.2025', duration: '22:05' },
    ];
    const startRecording = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            mediaRecorder.current = new MediaRecorder(stream);

            mediaRecorder.current.ondataavailable = (event) => {
                audioChunks.current.push(event.data);
            };

            mediaRecorder.current.onstop = () => {
                const blob = new Blob(audioChunks.current, { type: 'audio/wav' });
                setAudioBlob(blob);
                audioChunks.current = [];
            };

            mediaRecorder.current.start();
            setIsRecording(true);
            setAudioBlob(null); // Önceki kaydı temizle
        } catch (err) {
            console.error("Mikrofon erişim hatası:", err);
            setError("Mikrofon erişimi reddedildi. Lütfen tarayıcı ayarlarından izin verin.");
        }
    };

    const stopRecording = () => {
        if (mediaRecorder.current) {
            mediaRecorder.current.stop();
            setIsRecording(false);
        }
    };

    // --- METNE DÖNÜŞTÜRME FONKSİYONU ---
    const handleTranscribe = async () => {
        const fileToTranscribe = uploadedFile || audioBlob;
        if (!fileToTranscribe) {
            setError("Lütfen bir ses dosyası yükleyin veya bir ses kaydedin.");
            return;
        }
        setError('');
        setTranscription('');
        setIsLoading(true);

        try {
            // Artık token göndermemize gerek yok, authService hallediyor.
            const response = await transcribeAudio(fileToTranscribe);
            setTranscription(response.data.transcription);
        } catch (err) {
            if (err.response && err.response.status === 401) {
                setError("Oturum süreniz dolmuş veya geçersiz. Lütfen tekrar giriş yapın.");
                // İsteğe bağlı: localStorage.clear(); navigate('/login');
            } else {
                setError(err.response?.data?.error || "Metne dönüştürme sırasında bir hata oluştu.");
            }
        } finally {
            setIsLoading(false);
            setUploadedFile(null);
            setAudioBlob(null);
        }
    };
    const handleLogout = () => {
        localStorage.removeItem('authToken'); // Token'ı sil
        navigate('/login'); // Login sayfasına yönlendir
    };

    const renderContent = () => {
        switch (activeView) {
            case 'settings':
                return (
                    <div className="settings-view">
                        <div className="view-header">
                            <h3>Ayarlar</h3>
                            <button className="back-to-dash-button" onClick={() => setActiveView('main')}>
                                <BsArrowLeft /> Ana Sayfaya Dön
                            </button>
                        </div>
                        <div className="settings-content">
                            <div className="settings-form-container">
                                <h4>Kullanıcı Bilgileri</h4>
                                <form>
                                    <Input name="firstName" type="text" placeholder="Ad" defaultValue="Kullanıcı" />
                                    <Input name="lastName" type="text" placeholder="Soyad" defaultValue="Adı" />
                                    <Input name="email" type="email" placeholder="E-posta" defaultValue="kullanici@mail.com" disabled />
                                    <Button type="submit">Bilgileri Güncelle</Button>
                                </form>
                            </div>
                            <div className="settings-form-container">
                                <h4>Şifre Değiştir</h4>
                                <form>
                                    <Input name="currentPassword" type="password" placeholder="Mevcut Şifre" />
                                    <Input name="newPassword" type="password" placeholder="Yeni Şifre" />
                                    <Input name="confirmPassword" type="password" placeholder="Yeni Şifre (Tekrar)" />
                                    <Button type="submit">Şifreyi Değiştir</Button>
                                </form>
                            </div>
                        </div>
                    </div>
                );
            case 'history':
                return (
                    <div className="history-view">
                        <div className="view-header">
                            <h3>Geçmiş Transkriptler</h3>
                            <button className="back-to-dash-button" onClick={() => setActiveView('main')}>
                                <BsArrowLeft /> Ana Sayfaya Dön
                            </button>
                        </div>
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
                );
            default: // 'main'
                return (
                    <div className="dashboard-content">
                        <div className="recorder-section">
                            <h2>Ses Kaydet</h2>
                            <div className="recorder-visualizer">
                                {isRecording ? <BsRecordCircle className="mic-icon recording" /> : <BsFillMicFill className="mic-icon" />}
                                <div className="record-time">{audioBlob ? "Kayıt Tamamlandı" : (isRecording ? "Kaydediliyor..." : "00:00")}</div>
                            </div>
                            {isRecording ? (
                                <button className="record-button stop" onClick={stopRecording}><BsStopCircle /> Kaydı Durdur</button>
                            ) : (
                                <button className="record-button" onClick={startRecording}><BsFillMicFill /> Kaydı Başlat</button>
                            )}
                        </div>
                        <div className="upload-section">
                            <h2>Veya Bir Dosya Yükle</h2>
                            <div className="upload-area" onClick={() => document.querySelector('.file-input-hidden').click()}>
                                <BsUpload className="upload-icon" />
                                <p>{uploadedFile ? uploadedFile.name : "Dosyanızı buraya sürükleyin veya seçmek için tıklayın."}</p>
                                <input type="file" className="file-input-hidden" onChange={(e) => setUploadedFile(e.target.files[0])} style={{ display: 'none' }} />
                            </div>
                            <button className="transcribe-button" onClick={handleTranscribe} disabled={isLoading || (!audioBlob && !uploadedFile)}>
                                {isLoading ? 'Dönüştürülüyor...' : 'Metne Dönüştür'}
                            </button>
                        </div>
                    </div>
                );
        }
    };

    return (
        <div className="dashboard-page">
            <div className="background-icons">
                <div className="background-icon-wrapper icon1"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon2"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon3"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon4"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon5"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon6"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon7"><BsFillMicFill /></div>
            </div>

            <div className="dashboard-box">
                <header className="dashboard-header">
                    <h1>SpeeDh'e hoş geldin, {userName}!</h1>
                    <div className="header-buttons">
                        <button className="history-button" onClick={() => setActiveView('history')}>
                            <BsClockHistory /> Geçmiş Transkriptler
                        </button>
                        <button className="settings-button" onClick={() => setActiveView('settings')} title="Ayarlar">
                            <BsGearFill />
                        </button>
                        <div className="settings-form-container logout-section">
                            <Button onClick={handleLogout}>Çıkış Yap</Button>
                        </div>
                    </div>
                </header>

                <main className="dashboard-main-content">
                    {renderContent()}
                </main>
                {/* Hata ve Sonuç Alanı */}
                {(error || transcription || isLoading) && (
                    <footer className="dashboard-footer">
                        {isLoading && <p className="loading-message">Lütfen bekleyin, metin oluşturuluyor...</p>}
                        {error && <p className="error-message">{error}</p>}
                        {transcription && (
                            <div className="transcription-result">
                                <h4>Metin Çıktısı:</h4>
                                <p>{transcription}</p>
                            </div>
                        )}
                    </footer>
                )}
            </div>
        </div>
    );
};

export default DashboardPage;