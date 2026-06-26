#ifndef TELEMETRY_GENERATOR_H
#define TELEMETRY_GENERATOR_H

#include <string>
#include <map>
#include "telemetry.pb.h"

class TelemetryGenerator {
public:
    static telemetry::TelemetryPayload generatePayload(const std::string& agent_id);
private:
    static std::map<std::string, int> agent_uptimes;
};

#endif
