cmake_minimum_required(VERSION 3.10)
project(spm_agent)

set(CMAKE_CXX_STANDARD 17)

find_package(PkgConfig REQUIRED)
pkg_check_modules(PROTOBUF REQUIRED protobuf)
pkg_check_modules(GRPC REQUIRED grpc++)

set(PROTO_SRC "${CMAKE_CURRENT_BINARY_DIR}/telemetry.pb.cc")
set(PROTO_HEADER "${CMAKE_CURRENT_BINARY_DIR}/telemetry.pb.h")
set(GRPC_SRC "${CMAKE_CURRENT_BINARY_DIR}/telemetry.grpc.pb.cc")
set(GRPC_HEADER "${CMAKE_CURRENT_BINARY_DIR}/telemetry.grpc.pb.h")

add_custom_command(
    OUTPUT "${PROTO_SRC}" "${PROTO_HEADER}" "${GRPC_SRC}" "${GRPC_HEADER}"
    COMMAND protoc
    ARGS --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
         --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
         --plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin
         -I "${CMAKE_SOURCE_DIR}/proto"
         "${CMAKE_SOURCE_DIR}/proto/telemetry.proto"
    DEPENDS "${CMAKE_SOURCE_DIR}/proto/telemetry.proto"
)

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${PROTOBUF_INCLUDE_DIRS})
include_directories(${GRPC_INCLUDE_DIRS})

add_executable(spm_agent 
    src/main.cpp 
    src/telemetry_generator.cpp 
    src/api_client.cpp
    ${PROTO_SRC}
    ${GRPC_SRC}
)

target_link_libraries(spm_agent
    ${PROTOBUF_LIBRARIES}
    ${GRPC_LIBRARIES}
    pthread
)
