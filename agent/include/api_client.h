#ifndef API_CLIENT_H
#define API_CLIENT_H

#include <string>

class ApiClient {
public:
    ApiClient(const std::string& endpoint);
    void sendTelemetry(const std::string& payload);

private:
    std::string m_endpoint;
};

#endif // API_CLIENT_H
