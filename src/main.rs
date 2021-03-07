use snowflake_service;
use std::env;

use tonic::{transport::Server, Request, Response, Status};

use snowflake::snowflake_server::{Snowflake, SnowflakeServer};
use snowflake::{Empty, SnowflakeReply};

pub mod snowflake {
    tonic::include_proto!("snowflake");
}

#[derive(Debug, Default)]
pub struct MySnowflake {
    generator: snowflake_service::Generator,
}

impl MySnowflake {
    fn new(worker_id: u64) -> MySnowflake {
        let s = snowflake_service::Generator::new(worker_id);
        MySnowflake{
            generator: s,
        }
    }

}

#[tonic::async_trait]
impl Snowflake for MySnowflake {
    async fn get_snowflake(&self, _request: Request<Empty>) -> Result<Response<SnowflakeReply>, Status> {
        let id = match self.generator.generate_snowflake() {
            Ok(i) => i,
            Err(s) => return Err(Status::resource_exhausted(s)),
        };
        let reply = SnowflakeReply {
            id,
            id_str: id.to_string(),
        };
        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse()?;
    let worker_id: u64 = match env::var("WORKER_ID") {
        Ok(id) => id.parse::<u64>().unwrap(),
        Err(_) => panic!("No WORKER_ID env variable found"),
    };
    let s = MySnowflake::new(worker_id);
    
    Server::builder()
        .add_service(SnowflakeServer::new(s))
        .serve(addr)
        .await?;

    Ok(())
}
