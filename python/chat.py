import json
import httpx
import sys

def simple_chat(model="model", host="url"):
    print(f"Starting chat with {model}. Type 'exit' to end the conversation.")
    messages = []
    while True:
        user_input = input("\nYou: ")
        if user_input.lower() in ["exit", "quit", "bye"]:
            print("Goodbye!")
            break
        messages.append({
            "role": "user",
            "content": user_input
        })
        payload = {
            "model": model,
            "messages": messages,
            "stream": True
        }
        client = httpx.Client()
        assistant_response = ""
        print("\nAssistant: ", end="", flush=True)
        with client.stream("POST", f"{host}/", json=payload) as response:
            if response.status_code != 200:
                print(f"Error: {response.text}")
                continue
            for line in response.iter_lines():
                if line:
                    chunk = json.loads(line)
                    content_chunk = chunk.get("message", {}).get("content", "")
                    assistant_response += content_chunk
                    print(content_chunk, end="", flush=True)
        print()
        messages.append({
            "role": "assistant",
            "content": assistant_response
        })

if __name__ == "__main__":
    simple_chat()