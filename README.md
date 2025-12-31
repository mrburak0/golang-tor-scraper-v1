# web-scrapper-tor-browser-v1
Bu proje, Siber Vatan Programı Yıldız CTI Takımı görevleri kapsamında, eğitim ve öğrenme amacıyla geliştirilmiş Tor ağı üzerinden çalışan bir Web Scraper uygulamasıdır. Çalışma kapsamında Go (Golang) dili kullanılarak hedef web sitesinin URL’si komut satırından alınmış, istekler Tor Browser / Tor SOCKS5 proxy üzerinden yönlendirilmiş, hedef sayfanın statik HTML içeriği anonim olarak çekilmiş, sayfanın ekran görüntüsü alınmış ve elde edilen veriler yerel dosyalara kaydedilmiştir.

Proje, Siber Tehdit İstihbaratı (CTI) alanında, özellikle anonim veri toplama, kaynak gizleme ve Tor tabanlı açık kaynak istihbaratı (OSINT/CTI) süreçlerini uygulamalı olarak öğrenmeyi hedeflemektedir. Geliştirilen yapı, ilerleyen aşamalarda daha gelişmiş veri analizleri, farklı çıktı formatları, karanlık ağ (dark web) kaynaklarıyla uyumluluk ve ek CTI modülleri eklenebilecek şekilde genişletilmeye açıktır.

Kullanılan Teknolojiler
	•	Go (Golang)
	•	chromedp (Tor proxy üzerinden headless tarayıcı kontrolü)
	•	Tor / Tor Browser (SOCKS5 Proxy)

Oluşturulan Çıktılar
	•	Hedef web sayfasına ait HTML kaynak kodu dosyası
	•	Sayfa üzerinde bulunan bağlantıların listelendiği metin dosyası
	•	Hedef web sayfasına ait ekran görüntüsü (screenshot)

Sorumluluk Reddi

Bu proje tamamen eğitim amaçlı olarak geliştirilmiştir. Tor ağı üzerinden erişilen ve çekilen içeriklerin tüm fikrî mülkiyet hakları ilgili web sitelerine aittir. Proje; yetkisiz erişim, hukuka aykırı faaliyetler, mahremiyet ihlali veya kötü niyetli kullanım amacıyla tasarlanmamış olup, bu tür kullanımlar kesinlikle önerilmez ve sorumluluk kabul edilmez.
