use std::{
    io::{self, Read, Write},
    net::{Shutdown, TcpListener, TcpStream},
    process::exit,
};

fn main() -> io::Result<()> {
    let socket = TcpListener::bind("127.0.0.1:5555")?;

    for connection in socket.incoming() {
        match connection {
            Ok(stream) => std::thread::spawn(|| {
                handle_connection(stream).expect("handling connection failed");
            }),

            Err(err) => {
                eprintln!("error is {err}");
                exit(-1);
            }
        };
    }

    Ok(())
}

fn handle_connection(mut client_socket: TcpStream) -> io::Result<()> {
    let mut buf = [0; 1024];
    let received = client_socket.read(&mut buf)?;
    println!(
        "read {received} characters.\n\n{:?}\n",
        String::from_utf8_lossy(&buf[..received])
    );

    client_socket.write_all("hello!".as_bytes())?;

    client_socket.shutdown(Shutdown::Both)?;

    Ok(())
}
