from googletrans import Translator

translator = Translator(service_urls = ['translate.google.com'])
print("1. O'zbek tili")
print("2. Ingliz tili")
print("3. Rus tili")
print("4. Turk tili")
til = input("Tilni kiriting: ")
if til == 1:
    til = 'uz'
elif til == 2:
    til = 'en'
elif til == 3:
    til = 'ru'
elif til == 4:
    til = 'tr'
text = input("So'zni kiriting: ")
tarjima = translator.translate(text, dest = til)

print("natija:", tarjima.text)
# tr ru en uz