# Gerekli kütüphaneleri içe aktaralım
import noisereduce as nr
from pydub import AudioSegment
from scipy.io import wavfile
import numpy as np
import os # Dosya yolu işlemleri için

def clean_audio(input_file_path, output_file_path):
    """
    Bir ses dosyasındaki gürültüyü temizler ve sonucu kaydeder.

    Args:
        input_file_path (str): Gürültülü ses dosyasının yolu (örn: "ses.mp3").
        output_file_path (str): Temizlenmiş ses dosyasının kaydedileceği yol (örn: "temiz_ses.wav").
    """
    try:
        # 1. Ses dosyasını yükle (pydub sayesinde birçok formatı destekler)
        print(f"'{input_file_path}' dosyası yükleniyor...")
        audio = AudioSegment.from_file(input_file_path)

        # 2. Ses verisini NumPy dizisine çevir
        # Whisper ve noisereduce'un çalışması için bu format gereklidir.
        samples = np.array(audio.get_array_of_samples())
        
        # Eğer stereo ise mono'ya çevirelim (genellikle transkripsiyon için mono yeterlidir)
        if audio.channels > 1:
            samples = samples.reshape((-1, audio.channels)).mean(axis=1)

        print("Gürültü azaltma işlemi uygulanıyor...")
        # 3. noisereduce ile gürültüyü azalt
        # y: ses verisi (NumPy dizisi)
        # sr: örnekleme oranı (sample rate)
        reduced_noise_samples = nr.reduce_noise(y=samples.astype(np.float32), sr=audio.frame_rate)

        # 4. Temizlenmiş veriyi .wav formatında kaydet
        print(f"Temizlenmiş ses '{output_file_path}' olarak kaydediliyor...")
        wavfile.write(output_file_path, audio.frame_rate, reduced_noise_samples.astype(np.int16))
        
        print("İşlem başarıyla tamamlandı!")

    except FileNotFoundError:
        print(f"HATA: '{input_file_path}' dosyası bulunamadı. Lütfen dosya yolunu kontrol edin.")
    except Exception as e:
        print(f"Bir hata oluştu: {e}")

# --- KODU BURADA KULLANABİLİRSİNİZ ---

if __name__ == "__main__":
    # Temizlemek istediğiniz ses dosyasının adını ve yolunu belirtin
    # Bu script ile aynı klasörde olduğunu varsayıyoruz.
    input_audio = "ornek_gurultulu_ses.mp3" 
    
    # Çıktı dosyasının adını belirleyin. .wav formatında olmalı.
    output_audio = "temizlenmis_ses.wav"

    # Fonksiyonu çağırarak işlemi başlatın
    clean_audio(input_audio, output_audio)
    
    # Bu adımdan sonra "temizlenmis_ses.wav" dosyasını
    # Whisper AI modelinize vererek transkripsiyonu başlatabilirsiniz.