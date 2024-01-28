
#ifdef __MACH__

#include <CoreFoundation/CoreFoundation.h>
#include <CFNetwork/CFNetwork.h>

int main() {
    // 设置代理信息
    CFStringRef proxyHost = CFStringCreateWithCString(NULL, "proxy.example.com", kCFStringEncodingUTF8);
    CFNumberRef proxyPort = CFNumberCreate(NULL, kCFNumberIntType, &8080);

    // 创建代理字典
    CFDictionaryRef proxyDictKeys[] = { kCFNetworkProxiesHTTPEnable, kCFNetworkProxiesHTTPProxy, kCFNetworkProxiesHTTPPort };
    CFTypeRef proxyDictValues[] = { kCFBooleanTrue, proxyHost, proxyPort };
    CFDictionaryRef proxyDict = CFDictionaryCreate(NULL, (const void **)proxyDictKeys, (const void **)proxyDictValues, 3, &kCFTypeDictionaryKeyCallBacks, &kCFTypeDictionaryValueCallBacks);

    // 设置代理
    CFStringRef proxyKey = CFStringCreateWithCString(NULL, "HTTPProxy", kCFStringEncodingUTF8);
    CFNetworkSetSystemProxySettings(proxyKey, proxyDict);

    // 释放资源
    CFRelease(proxyHost);
    CFRelease(proxyPort);
    CFRelease(proxyDict);
    CFRelease(proxyKey);

    return 0;
}
#endif