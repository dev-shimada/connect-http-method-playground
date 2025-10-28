import asyncio
import httpx
from gen.proto.api.v1.api_connect import ApiServiceClient
from gen.proto.api.v1 import api_pb2


async def main():
    base_url = "http://api:8081"

    async with httpx.AsyncClient() as session:
        client = ApiServiceClient(address=base_url, session=session)

        # Call Post RPC
        print("Calling Post RPC...")
        post_request = api_pb2.PostRequest(
            user_id="user123",
            text="Hello from Python Connect RPC client!"
        )

        post_response = await client.post(post_request)
        print(f"Post Response - ID: {post_response.id}")

        # Call Get RPC with the returned ID
        print(f"\nCalling Get RPC with ID: {post_response.id}...")
        get_request = api_pb2.GetRequest(id=post_response.id)

        get_response = await client.get(get_request)
        print(f"Get Response - User ID: {get_response.user_id}, Text: {get_response.text}")



if __name__ == "__main__":
    asyncio.run(main())
