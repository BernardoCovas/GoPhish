import os
import requests

FILENAME = "index.raw.html"
INDEX = "index.html"
FOLDER = "__res__"

if not os.path.exists(FOLDER):
    os.makedirs(FOLDER)

with open(FILENAME) as f:
    with open(INDEX, "w") as of:
        c = 0
        for _i, line in enumerate(f):

            if "https://" in line:
                c += 1
                i = line.find("=")
                web = line[i+2:-3]

                if web == "https://www.facebook.com/":
                    continue
 
                ext = web.split(".")[-1]

                res = requests.get(web, stream=True)
                filepath = os.path.join(FOLDER, str(_i)) + "." + ext
                
                with open(filepath, "wb") as f:
                    for chunk in res.iter_content(2048):
                        f.write(chunk)

                if "https://m.facebook.com/login/" in web:
                    filepath = "login/"
                 
                line = line.replace(web, f"/{filepath}")
                print(web)
                print(filepath)
                print(line)

            of.write(line)
