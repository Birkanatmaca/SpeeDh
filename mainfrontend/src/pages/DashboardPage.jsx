import React, { useState, useRef } from 'react';
import { BsFillMicFill, BsUpload, BsClockHistory, BsGearFill, BsArrowLeft, BsRecordCircle, BsStopCircle, BsTrashFill, BsBoxArrowRight, BsEyeFill } from "react-icons/bs";
import Input from '../components/common/Input';
import Button from '../components/common/Button';
import { transcribeAudio, getHistory, getAudioFile, deleteTranscript } from '../api/authService';
import { useNavigate } from 'react-router-dom';
import './DashboardPage.scss';
import BiacaButton from '../components/common/BiacaButton';
import Modal from '../components/common/Modal';
import ConfirmationModal from '../components/common/ConfirmationModal';
const DashboardPage = () => {
    const [activeView, setActiveView] = useState('main');
    const [isRecording, setIsRecording] = useState(false);
    const [audioBlob, setAudioBlob] = useState(null);
    const [audioURL, setAudioURL] = useState('');
    const [uploadedFile, setUploadedFile] = useState(null);
    const [transcription, setTranscription] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const [transcripts, setTranscripts] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedTranscript, setSelectedTranscript] = useState(null);

    const mediaRecorder = useRef(null);
    const audioChunks = useRef([]);

    const [isConfirmModalOpen, setIsConfirmModalOpen] = useState(false);
    const [transcriptToDelete, setTranscriptToDelete] = useState(null);

    const userName = "Kullanıcı";
    const navigate = useNavigate();

    const startRecording = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            mediaRecorder.current = new MediaRecorder(stream);
            mediaRecorder.current.ondataavailable = (event) => {
                audioChunks.current.push(event.data);
            };
            mediaRecorder.current.onstop = () => {
                // ---------- DÜZELTME 1: Blob formatını tekrar 'audio/wav' yapıyoruz ----------
                const blob = new Blob(audioChunks.current, { type: 'audio/wav' });
                const url = URL.createObjectURL(blob);
                setAudioBlob(blob);
                setAudioURL(url);
                audioChunks.current = [];
            };
            mediaRecorder.current.start();
            setIsRecording(true);
            handleDeleteRecording();
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

    const handleDeleteRecording = () => {
        setAudioBlob(null);
        setAudioURL('');
    };
    const handleDeleteClick = (transcript) => {
        setTranscriptToDelete(transcript);
        setIsConfirmModalOpen(true);
    };

    // Onay modalını kapatır
    const handleCloseConfirmModal = () => {
        setIsConfirmModalOpen(false);
        setTranscriptToDelete(null);
    };

    // Onay modalında "Evet"e basıldığında silme işlemini gerçekleştirir
    const handleConfirmDelete = async () => {
        if (!transcriptToDelete) return;

        try {
            await deleteTranscript(transcriptToDelete.id);
            // Başarılı olursa, listeden silinmiş olanı filtrele
            setTranscripts(transcripts.filter(t => t.id !== transcriptToDelete.id));
        } catch (err) {
            console.error("Silme hatası:", err);
            setError("Kayıt silinirken bir hata oluştu.");
        } finally {
            // Modalı kapat
            handleCloseConfirmModal();
        }
    };

    const handleTranscribe = async () => {
        let fileToTranscribe;

        if (uploadedFile) {
            fileToTranscribe = uploadedFile;
        } else if (audioBlob) {
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, '0');
            const day = String(now.getDate()).padStart(2, '0');
            const hours = String(now.getHours()).padStart(2, '0');
            const minutes = String(now.getMinutes()).padStart(2, '0');
            const seconds = String(now.getSeconds()).padStart(2, '0');

            // ---------- DÜZELTME 2: Dosya adını ve formatını tekrar '.wav' yapıyoruz ----------
            const fileName = `Kayıt_${year}${month}${day}_${hours}${minutes}${seconds}.wav`;
            fileToTranscribe = new File([audioBlob], fileName, { type: "audio/wav" });
        }

        if (!fileToTranscribe) {
            setError("Lütfen bir ses dosyası yükleyin veya bir ses kaydedin.");
            return;
        }

        setError('');
        setTranscription('');
        setIsLoading(true);

        try {
            const response = await transcribeAudio(fileToTranscribe);
            setTranscription(response.data.transcription);
        } catch (err) {
            console.error("API Hatası:", err.response || err);
            if (err.response && err.response.status === 401) {
                setError("Oturum süreniz dolmuş veya geçersiz. Lütfen tekrar giriş yapın.");
            } else {
                setError(err.response?.data?.error || "Metne dönüştürme sırasında bir hata oluştu.");
            }
        } finally {
            setIsLoading(false);
            setUploadedFile(null);
            handleDeleteRecording();
        }
    };

    const handleHistoryClick = async () => {
        setActiveView('history');
        setIsLoading(true);
        setError('');
        try {
            const response = await getHistory();
            setTranscripts(response.data || []);
        } catch (err) {
            console.error("Geçmiş yüklenirken hata:", err.response || err);
            const errorMessage = err.response?.data?.error || "Geçmiş transkriptler yüklenemedi. Sunucuya ulaşılamıyor veya bir hata oluştu. Lütfen tekrar deneyin.";
            setError(errorMessage);
            setTranscripts([]);
        } finally {
            setIsLoading(false);
        }
    };

    const formatDate = (dateString) => {
        const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
        return new Date(dateString).toLocaleDateString('tr-TR', options);
    };

    const handleLogout = () => {
        localStorage.removeItem('authToken');
        navigate('/login');
    };

    const handleViewClick = (transcript) => {
        setSelectedTranscript(transcript);
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setSelectedTranscript(null);
    };

    const handleDownloadText = () => {
        if (!selectedTranscript) return;
        const blob = new Blob([selectedTranscript.transcription_text], { type: 'text/plain;charset=utf-8' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${selectedTranscript.title}.txt`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    };

    const handleDownloadAudio = async () => {
        if (!selectedTranscript) return;
        try {
            const response = await getAudioFile(selectedTranscript.id);
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const a = document.createElement('a');
            a.href = url;
            a.download = selectedTranscript.original_filename || 'audio.wav';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (err) {
            console.error("Ses dosyası indirilemedi:", err);
            setError("Ses dosyası indirilirken bir hata oluştu.");
        }
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
                            {isLoading ? <p>Yükleniyor...</p> : error ? <div className="error-message" style={{ textAlign: 'center', padding: '1rem' }}>{error}</div> : (
                                <table className="history-table">
                                    <thead>
                                        <tr>
                                            <th>Başlık</th>
                                            <th>Tarih</th>
                                            <th>İşlemler</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {transcripts.length > 0 ? transcripts.map(item => (
                                            <tr key={item.id}>
                                                <td>{item.title}</td>
                                                <td>{formatDate(item.created_at)}</td>
                                                <td>
                                                    <button className="action-button view" onClick={() => handleViewClick(item)}><BsEyeFill /> Görüntüle</button>
                                                    {/* OnClick olayını güncelle */}
                                                    <button className="action-button delete" onClick={() => handleDeleteClick(item)}><BsTrashFill /> Sil</button>
                                                </td>
                                            </tr>
                                        )) : (
                                            <tr>
                                                <td colSpan="3" style={{ textAlign: 'center', padding: '1rem' }}>Henüz bir transkript oluşturmadınız.</td>
                                            </tr>
                                        )}
                                    </tbody>
                                </table>
                            )}
                        </main>
                    </div>
                );
            default: // 'main'
                return (
                    <div className="transcribe-container">
                        <div className="input-options">
                            <div className="recorder-section">
                                <h2>Ses Kaydet</h2>
                                <div className="recorder-visualizer">
                                    {audioURL ? (
                                        <div className="audio-player-container">
                                            <audio src={audioURL} controls />
                                            <button className="delete-button" onClick={handleDeleteRecording} title="Kaydı Sil">
                                                <BsTrashFill />
                                            </button>
                                        </div>
                                    ) : (
                                        <>
                                            {isRecording ? <BsRecordCircle className="mic-icon recording" /> : <BsFillMicFill className="mic-icon" />}
                                            <div className="record-time">{isRecording ? "Kaydediliyor..." : "00:00"}</div>
                                        </>
                                    )}
                                </div>
                                {!audioBlob && (isRecording ? (
                                    <button className="record-button stop" onClick={stopRecording}><BsStopCircle /> Kaydı Durdur</button>
                                ) : (
                                    <button className="record-button" onClick={startRecording}><BsFillMicFill /> Kaydı Başlat</button>
                                ))}
                            </div>
                            <div className="upload-section">
                                <h2>Veya Bir Dosya Yükle</h2>
                                <div className="upload-area" onClick={() => document.querySelector('.file-input-hidden').click()}>
                                    <BsUpload className="upload-icon" />
                                    <p>{uploadedFile ? uploadedFile.name : "Dosyanızı buraya sürükleyin veya seçmek için tıklayın."}</p>
                                    <input type="file" className="file-input-hidden" onChange={(e) => setUploadedFile(e.target.files[0])} style={{ display: 'none' }} accept="audio/*" />
                                </div>
                            </div>
                        </div>
                        <div className="transcribe-action">
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
            <BiacaButton />
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
                    <h1>SpeeDch'e hoş geldin, {userName}!</h1>
                    <div className="header-buttons">
                        <button className="history-button" onClick={handleHistoryClick}>
                            <BsClockHistory /> Geçmiş Transkriptler
                        </button>
                        <button className="settings-button" onClick={() => setActiveView('settings')} title="Ayarlar">
                            <BsGearFill />
                        </button>
                        <button className="logout-button" onClick={handleLogout} title="Çıkış Yap">
                            <BsBoxArrowRight />
                        </button>
                    </div>
                </header>

                <main className="dashboard-main-content">
                    {renderContent()}
                </main>
                {(error && activeView === 'main' || transcription || (isLoading && activeView === 'main')) && (
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

            <Modal
                isOpen={isModalOpen}
                onClose={handleCloseModal}
                transcript={selectedTranscript}
                onDownloadText={handleDownloadText}
                onDownloadAudio={handleDownloadAudio}
            />
            <ConfirmationModal
                isOpen={isConfirmModalOpen}
                onClose={handleCloseConfirmModal}
                onConfirm={handleConfirmDelete}
                message="Kaydı ve transkripti silmek istediğinize emin misiniz?"
            />
        </div>
    );
};

export default DashboardPage;