use reqwest::{Error};
use std::result::Result;
use crate::commands::client::AppClient;

pub async fn rm(ac: &mut impl AppClient, app_name: &str)
-> Result<(), Error> {
    ac.rm_app(app_name).await
}

#[cfg(test)]
mod tests {
    use super::*;
    use futures::executor::block_on;
    use crate::commands::client::TestAppClient;


    #[test]
    fn test_run() {
        let mut cl = TestAppClient::new();
        let res_fut = rm(&mut cl, "testapp");
        let res = block_on(res_fut).unwrap();
        
        assert_eq!(res, ());
        assert_eq!(cl.rm_counter, 1);
        assert_eq!(cl.add_counter, 0);
    }

}
