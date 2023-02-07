var MianProxy = "PROXY localhost:8080";

var prolicy = [];

function addRegDom(base, area, sub){
  return prolicy.push(RegExp(
      ((sub && "\\." + sub ) || "") +
      base.replace(".","\\.") +
      ((area && "\\." + area) || "") + "$"
    ));
};

addRegDom("google", "((com\\.)?[a-z]{2}|com)?");
addRegDom("youtube", "((com\\.)?[a-z]{2}|com)?");
addRegDom("blogspot.com");
addRegDom("pixiv.net");
addRegDom("steamcommunity.com");
addRegDom("twitter.com");

function FindProxyForURL(url, host) {
  for (var i=0; i < prolicy.length; i++){
      if (prolicy[i] instanceof RegExp && prolicy[i].test(host)) {
          return MianProxy;
      }
  }
  alert("non-regular: "+url)
  return (() => {
    if (host == "sockstest.com") {
      return "SOCKS mysocks.localhost:11080; SOCKS5 mysocks.localhost; SOCKS4 mysocks.localhost"
    }
    if (host == "complex.com") {
      return "HTTP a.localhost; HTTPS b.localhost; SOCKS4 c.localhost; NATIVE"
    }
    if (host == "error.com") {
      return "HTTP a.localhost; XPROTO b.localhost; NOTHING"
    }
    if (host == "nothing.com") {
      return ""
    }
    return 'DIRECT';
  })()
}
