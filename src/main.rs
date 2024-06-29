use std::{
    io,
    net::{IpAddr, Ipv4Addr, SocketAddr, UdpSocket},
};

use nix::sys::socket::{listen, Backlog};

fn main() -> Result<(), io::Error> {
    let socket_addr = SocketAddr::new(IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)), 8080);

    let socket = UdpSocket::bind(&socket_addr)?;
    let backlog = Backlog::new(1)?;

    listen(&socket, backlog)?;

    let mut buf = [0; 1024];
    loop {
        match socket.recv(&mut buf) {
            Ok(received) => println!("{} recieved.\nmesseage:\n{:?}", received, &buf[..received]),
            Err(_) => break,
        }
    }

    Ok(())
}
