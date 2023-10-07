import argparse
import time

def main():
    parser = argparse.ArgumentParser(description="Sleep for a specified amount of time.")
    
    # Add the -t argument
    parser.add_argument('-t', type=int, help='Number of seconds to sleep', required=True)
    
    # Parse the arguments
    args = parser.parse_args()
    
    # Sleep for the specified amount of time
    print(f"Sleeping for {args.t} seconds...")
    time.sleep(args.t)
    print("Done sleeping!")

if __name__ == "__main__":
    main()