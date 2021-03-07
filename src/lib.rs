use std::time;
use std::sync::Mutex;
use std::convert::TryInto;

const EPOCH_TIME: time::Duration = time::Duration::from_millis(1609459200000);
const TIME_MASK: u64 = 0x1FFFFFFFFFF;

#[derive(Debug, Default)]
pub struct Generator {
    last_generated_time: Mutex<u64>,
    counter: Mutex<u64>,
    last_counter_rollover: Mutex<u64>,
    worker_id: u64,
}

impl Generator {
    pub fn new(worker_id: u64) -> Generator {
        if worker_id > 1023 {
            panic!("Invalid worker id")
        }
        Generator {
            last_generated_time: Mutex::new(0),
            counter: Mutex::new(0),
            last_counter_rollover: Mutex::new(0),
            worker_id: worker_id << 12,
        }
    }

    pub fn generate_snowflake(&self) -> Result<u64, &str> {
        let time: u64 = match time::SystemTime::now().duration_since(time::SystemTime::UNIX_EPOCH + EPOCH_TIME) {
            Ok(t) => t.as_millis().try_into().unwrap(),
            Err(_) => return Err("The epoch appears to be in the past! Check the system time."),
        };
        let mut last_generated_time = self.last_generated_time.lock().unwrap();
        let mut counter = self.counter.lock().unwrap();
        let mut last_counter_rollover = self.last_counter_rollover.lock().unwrap();
        if time < *last_generated_time {
            return Err("The current time is less than the last generated time! Check the system time.");
        }
        if *counter == 4095 {
            if *last_counter_rollover >= time {
                return Err("Too many requests in the current ms");
            }
            *last_counter_rollover = time;
            *counter = 0;
        } else {
            *counter += 1;
        }
        *last_generated_time = time;
        Ok(((time & TIME_MASK) << 22) + self.worker_id + *counter)
    }
}
