use reqwest::{Error};
use std::result::Result;
use crate::commands::client::AppClient;


pub async fn run(ac: &mut impl AppClient, app_name: &str, image: &str, port: u32) 
-> Result<(), Error> {
    ac.add_app(app_name, image, port).await
}

#[cfg(test)]
mod tests {
    use super::*;
    use futures::executor::block_on;
    use crate::commands::client::TestAppClient;


    #[test]
    fn test_run() {
        let mut cl = TestAppClient::new();
        let res_fut = run(&mut cl, "testapp", "testimage", 9090);
        let res = block_on(res_fut).unwrap();
        
        assert_eq!(res, ());
        assert_eq!(cl.rm_counter, 0);
        assert_eq!(cl.add_counter, 1);
    }

}
