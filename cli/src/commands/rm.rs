use reqwest::{Error};
use std::result::Result;
use crate::commands::client::AppClient;

pub async fn rm(ac: &mut impl AppClient, app_name: &str)
-> Result<(), Error> {
    ac.rm_app(app_name).await
}
