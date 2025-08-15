import { useState } from 'react';
import Input from '../components/common/Input';
import Button from '../components/common/Button';
import { login } from '../api/authService';
import { BsFillMicFill } from "react-icons/bs";
import './AuthPage.scss'; // Kayıt sayfasıyla aynı ortak stil dosyasını kullanıyoruz
import { useNavigate } from 'react-router-dom'; // Yönlendirme için eklendi

const LoginPage = () => {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate(); // Yönlendirme için eklendi

    // Input alanlarındaki her değişiklikte state'i güncelleyen fonksiyon
    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    // Form gönderildiğinde çalışacak fonksiyon
    const handleSubmit = async (e) => {
        e.preventDefault(); // Sayfanın yeniden yüklenmesini engelle
        setError('');
        setLoading(true); // Butonu devre dışı bırakıp "Giriş Yapılıyor..." yazısını göstermek için
        try {
            const response = await login(formData.email, formData.password);
            localStorage.setItem('authToken', response.data.token);

            // --- BAŞARILI GİRİŞ SONRASI ---
            console.log('Giriş Başarılı! Token:', response.data.token);
            // Dashboard sayfasına yönlendirme
            navigate('/dashboard');

        } catch (err) {
            // Backend'den gelen hata mesajını state'e kaydedip ekranda gösteriyoruz
            setError(err.response?.data?.error || 'Giriş sırasında bir hata oluştu.');
        } finally {
            // İşlem başarılı da olsa başarısız da olsa, loading durumunu bitir.
            setLoading(false);
        }
    };

    return (
        <div className="auth-page">
            {/* Arka Plan Animasyon Katmanı */}
            <div className="background-icons">
                <div className="background-icon-wrapper icon1"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon2"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon3"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon4"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon5"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon6"><BsFillMicFill /></div>
                <div className="background-icon-wrapper icon7"><BsFillMicFill /></div>
            </div>

            {/* Ana İçerik Konteyneri */}
            <div className="auth-container">

                {/* Sol Panel: Hareketli ve Renkli Kısım */}
                <div className="auth-panel">
                    <h1>Tekrar Hoş Geldin!</h1>
                    <p>Kaldığın yerden devam et. Fikirlerin seni bekliyor.</p>
                </div>

                {/* Sağ Panel: Formun Olduğu Kısım */}
                <div className="auth-form-wrapper">
                    <form className="auth-form" onSubmit={handleSubmit}>
                        <h2>Giriş Yap</h2>
                        <Input
                            name="email"
                            type="email"
                            placeholder="E-posta Adresi"
                            value={formData.email}
                            onChange={handleChange}
                        />
                        <Input
                            name="password"
                            type="password"
                            placeholder="Şifre"
                            value={formData.password}
                            onChange={handleChange}
                        />
                        <Button type="submit" disabled={loading}>
                            {loading ? 'Giriş Yapılıyor...' : 'Giriş Yap'}
                        </Button>

                        {error && <p className="error-message">{error}</p>}

                        <p className="form-footer">
                            Şifreni mi unuttun? <a href="/forgot-password">Sıfırla</a>
                        </p>
                        <p className="form-footer">
                            Hesabın yok mu? <a href="/register">Kayıt Ol</a>
                        </p>
                    </form>
                </div>

            </div>
        </div>
    );
};

export default LoginPage;