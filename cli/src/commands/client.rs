use reqwest::{Error, Client};
use std::result::Result;
use std::collections::HashMap;
use async_trait::async_trait;



#[async_trait]
pub trait AppClient {
    async fn add_app(&mut self, app_name: &str, app_image: &str, port: u32) -> Result<(), Error>;
    async fn rm_app(&mut self, app_name: &str) -> Result<(), Error>;
}

pub struct ProdAppClient {
    base_deploy_url: String
}

impl ProdAppClient {
    pub fn new(base_deploy_url: &str) -> ProdAppClient {
        ProdAppClient{
            base_deploy_url: base_deploy_url.to_string(),
        }
    }
}

#[async_trait]
impl AppClient for ProdAppClient {
    async fn add_app(&mut self, app_name: &str, app_image: &str, port: u32)
    -> Result<(), Error> {
        let port_string = port.to_string();
        let client = Client::new();
        let mut map = HashMap::new();
        map.insert("name", app_name);
        map.insert("image", &app_image);
        map.insert("port", &port_string);
    
        let request_url = format!("{}?name={}", self.base_deploy_url, app_name);
        let res_future = client.post(&request_url)
        .json(&map)
        .send();

        res_future.await.map(|_| ())
    }

    async fn rm_app(&mut self, app_name: &str)
    -> Result<(), Error> {
        let client = Client::new();
        let mut map = HashMap::new();
        map.insert("name", app_name);
        
        let request_url = format!("{}?name={}", self.base_deploy_url, app_name);
        let res_future = client.delete(&request_url)
        .json(&map)
        .send();
    
        res_future.await.map(|_| ())
    }
}

pub struct TestAppClient {
    pub add_counter: u32,
    pub rm_counter: u32,
    // add_return: Result<(), Error>,
    // rm_return: Result<(), Error>,
}

impl TestAppClient {
    // pub fn new(add_return: Result<(), Error>, rm_return: Result<(), Error>)
    // -> TestAppClient {
    //     TestAppClient{
    //         add_return: add_return,
    //         rm_return: rm_return,
    //     }
    // }
    pub fn new()
    ->TestAppClient {
        TestAppClient{add_counter: 0, rm_counter: 0}
    }
}

#[async_trait]
impl AppClient for TestAppClient {
    async fn add_app(&mut self, _: &str, _: &str, _: u32)
    -> Result<(), Error> {
        self.add_counter+=1;
        Ok(())
    }
    async fn rm_app(&mut self, _: &str)
    -> Result<(), Error> {
        self.rm_counter+=1;
        Ok(())
    }
}
