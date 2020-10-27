use reqwest::{Client, Response, Error};
use std::collections::HashMap;
use std::result::Result;

pub async fn rm(admin_url: &String, app_name: &String) 
-> Result<Response, Error> {
    println!("removing app {}", app_name);
    println!("URL: {}", admin_url);
    let client = Client::new();
    // TODO: maybe use Serde & a struct to serialize JSON
    let mut map = HashMap::new();
    map.insert("name", app_name);
    
    let req_url = format!("{}?name={}", admin_url, app_name);
    let res_future = client.delete(&req_url)
    .json(&map)
    .send();

    res_future.await
}
