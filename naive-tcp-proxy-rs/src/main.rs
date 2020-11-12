use futures::future::try_join;
use std::env;
use std::process;
use tokio::io;
use tokio::io::AsyncWriteExt;
use tokio::net::{TcpListener, TcpStream};

#[tokio::main]
async fn main() -> Result<(), io::Error> {
    if env::args().len() != 3 {
        println!(
            "usage: {} <listen-addr> <upstream-addr>",
            env::args().next().unwrap()
        );
        process::exit(1);
    }

    let listen_addr = env::args().nth(1).unwrap();
    let upstream_addr = env::args().nth(2).unwrap();

    let li = TcpListener::bind(listen_addr).await?;

    while let Ok((downstream, _)) = li.accept().await {
        tokio::spawn(proxy(downstream, upstream_addr.clone()));
    }

    Ok(())
}

async fn proxy(mut downstream: TcpStream, upstream_addr: String) {
    let mut upstream = TcpStream::connect(upstream_addr).await.unwrap();

    let (mut rd, mut wd) = downstream.split();
    let (mut ru, mut wu) = upstream.split();

    let client_to_server = async {
        io::copy(&mut rd, &mut wu).await?;
        wu.shutdown().await
    };

    let server_to_client = async {
        io::copy(&mut ru, &mut wd).await?;
        wd.shutdown().await
    };

    try_join(client_to_server, server_to_client).await;
}
