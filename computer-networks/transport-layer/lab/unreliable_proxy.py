# -*- coding: utf-8 -*-
import random
import socket
import sys


DROP_RATE = 0.6
CORRUPTION_RATE = 0.0


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Usage: python3 unreliable_proxy.py destination_port')
        sys.exit(1)

    dest = int(sys.argv[1])
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.bind(('', 0))
        host, port = sock.getsockname()
        print(f'Forwarding {host}:{port} -> 127.0.0.1:{dest}')

        while True:
            payload, address = sock.recvfrom(4096)
            print('\n⭐⭐ New Packet ⭐⭐\n')
            print(payload)

            if random.random() < DROP_RATE:
                print('\n🛑 PACKET DROPPED!\n')
                continue

            if random.random() < CORRUPTION_RATE:
                print('\n🔥 CORRUPTION!\n')
                payload = list(payload)
                payload[random.randrange(len(payload))] ^= 0xff
                payload = bytes(payload)
                # Send NACK
                sock.sendto(bytes.fromhex('01'), address)
            else:
                print('\n👌 OK\n')
                # Send ACK
                sock.sendto(bytes.fromhex('00'), address)
                # Forward to proxy
                sock.sendto(payload, ('', dest))


    finally:
        sock.close()
