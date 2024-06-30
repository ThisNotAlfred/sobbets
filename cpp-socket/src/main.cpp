#include <memory>
#include <print>
#include <string>
#include <vector>

#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>

auto
main(int argc, char* argv[]) -> int
{
    if (argc < 3) {
        std::print(stderr, "wrong number of argumets!\n");
        exit(-1);
    }

    // creting a udp6 socket
    auto hints        = addrinfo {};
    hints.ai_family   = AF_INET;
    hints.ai_socktype = SOCK_STREAM;

    // getting `addrinfo`s and putting them inside the `servinfos`
    auto servinfos = std::array<addrinfo*, 1> { nullptr };
    int result     = getaddrinfo(argv[1], argv[2], &hints, servinfos.data());
    if (result != 0) {
        std::print(stderr,
                   "I HATE MANUAL MEMORY MANAGEMENT. GIMME RUST BACK!. BTW THE ERROR WAS: {}\n",
                   gai_strerror(result));
        exit(-1);
    }

    // iterating through `servinfos` and sending a packet to them.
    // NOTE: WHAT THE HELL?! WHAT IS THIS ITERATION METHOD? IT'S GIVING VERY 1990s VIBES
    for (auto& servinfo = servinfos[0]; servinfo != nullptr; servinfo = servinfo->ai_next) {
        auto sockfd = socket(servinfo->ai_family, servinfo->ai_socktype, servinfo->ai_protocol);
        if (sockfd == -1) {
            std::print(stderr, "address info didn't work. trying next one...\n");
            continue;
        }

        // binding to the socket to read from it
        if (connect(sockfd, servinfo->ai_addr, servinfo->ai_addrlen) == -1) {
            std::print(stderr, "couldn't connect to the socket\n");
            close(sockfd);
            continue;
        }

        // sending packet "HELLO!" to each working socket
        auto sent = write(sockfd, "HELLO!\n", 7);
        if (sent == -1) {
            std::print(stderr, "failed to write to socket\n");
            close(sockfd);
            continue;
        }

        // pre-allocating buffer for reading from the socket
        auto buffer     = std::array<char, 1024> {};
        auto read_bytes = read(sockfd, buffer.data(), buffer.size());
        if (read_bytes == -1) {
            std::print(stderr, "reading from socket failed!\n");
            close(sockfd);
            continue;
        }

        // if reading didn't fail, we print the result.
        std::print("{}", buffer.data());
        // and close the socket at the end
        close(sockfd);
    }

    return 0;
}