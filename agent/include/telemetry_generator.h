#ifndef TELEMETRY_GENERATOR_H
#define TELEMETRY_GENERATOR_H

#include <string>

class TelemetryGenerator {
public:
    static std::string generatePayload(const std::string& agent_id);
};

#endif // TELEMETRY_GENERATOR_H
