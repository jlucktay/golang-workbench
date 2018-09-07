# secrets

**This tiny package acts as an extremely flimsy layer of security.**

Place your super secure secret into any given JSON file with a key named `token` then call `secrets.ReadTokenFromSecrets` with the path to the JSON file as an argument, and it will give you back the secret that is safely tucked away within the file.

Also, you should probably not check the secret JSON file(s) into git either!
