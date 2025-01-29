# LLM Benchmarking & Latency Testing Tool (Go)

🚀 A Golang-based tool for benchmarking LLM (Large Language Model) response times, measuring embedding performance, analyzing latency, and evaluating contextual reasoning accuracy.

## 🛠 Features
- **Multi-Model Support**: Runs benchmarks on various LLM models (e.g., StableLM2, LLaMA 3.1 70B Instruct).
- **Latency & Performance Testing**:
  - Measures **embedding time**, **response time**, and **first-byte latency**.
  - Provides a **structured table report** of all evaluations.
- **Cosine Similarity Check**: Uses **BGE-Base-en V1.5** to compare vector embeddings.
- **Contextual Reasoning Check**: Utilizes **LLaMA 3.1 70B Instruct** to assess reasoning capabilities.
- **Retrieval-Augmented Generation (RAG) Evaluation**: Supports document retrieval analysis.
- **Streaming Response Handling**: Captures real-time responses for precise benchmarking.

## 📥 Installation
1. **Clone the repository**:
   ```sh
   git clone https://github.com/RezaArani/llm-tester-go
   cd llm-tester-go
   ```

2. **Install dependencies**:
   ```sh
   go mod tidy
   ```

3. **Run the tool**:
   ```sh
   go run main.go
   ```
** To run this tool, you need Redis and Ollama. Installation instructions for both can be found online.

## 🔧 Configuration
Modify the `models` array in `main.go` to specify which models you want to test:
```go
 

	models := []string{
		"stablelm2",
		"llama3.2:1b",
		"llama3.1",
		"deepseek-r1:1.5b",
		"tinyllama",
		"llama3.2",
		"stablelm2",
		"aya",
		"granite3-dense:8b",
		"mistral",
		"wizardlm2",
		"qwen2",

		
	}

 
	queries := []string{
		"What is SemMapas? (limiting it to a maximum of 50 words)", // limit text for testing token count
		"Sum of SemMapas annual costs?",                            // me may keep model decide the result
		// "Say Hi",
	}
	
	const enRawText = `Welcome to SemMapas, your strategic partner in enhancing local engagement and tourism development. Designed specifically for businesses and municipalities, SemMapas offers a powerful platform to connect with residents and visitors alike, driving growth and prosperity in your community.
Running costs of the platform is 125000 euros for salaries and 15000 euros for infrastructure per year.
Our platform goes beyond traditional mapping services, offering a comprehensive suite of features tailored to meet the diverse needs of event organizers and businesses alike. From tourism guides to event navigation, SemMapas empowers you to create immersive experiences that captivate your audience and enhance their journey.
Our project has been launched since 2023 in Portugal.`

```
 

## 📊 Example Output
```

Results Table:
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Model                Query                                    ResponseTime    CONTEXT    PromptTokens     CompletionTokens TotalTokens      Cosine Similarity    Contextual Reasoning,Correct? 
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
stablelm2            Sum of SemMapas annual costs?            28.3218774s     1          259              150              409                0.7410108235066566   No,                           

stablelm2            Sum of SemMapas annual costs?            37.3171889s     0          74               313              387                0.6445678523262685   No,                           

deepseek-r1:1.5b     Sum of SemMapas annual costs?            1m4.7651514s    1          244              327              571                0.7335673185458016   Yes.                          

deepseek-r1:1.5b     Sum of SemMapas annual costs?            1m52.1123971s   0          59               828              887                0.661910460773835    No,                           


====================================================== Embedding Time =====================================================================
Model: stablelm2	 Embedding Time:6.4216s
Model: deepseek-r1:1.5b	 Embedding Time:14.4520387s
Context size: 770 bytes
Context: Welcome to SemMapas, your strategic partner in enhancing local engagement and tourism development. Designed specifically for businesses and municipalities, SemMapas offers a powerful platform to connect with residents and visitors alike, driving growth and prosperity in your community.
Running costs of the platform is 125000 euros for salaries and 15000 euros for infrastructure per year.
Our platform goes beyond traditional mapping services, offering a comprehensive suite of features tailored to meet the diverse needs of event organizers and businesses alike. From tourism guides to event navigation, SemMapas empowers you to create immersive experiences that captivate your audience and enhance their journey.
Our project has been launched since 2023 in Portugal.


====================================================== REPORT DETAILS =====================================================================

Model:stablelm2
Query:Sum of SemMapas annual costs?

-------------------------------------------------------------------------------------------------------------
               Vector Search       LLM Call            First Chunk         Stream Generation Time
-------------------------------------------------------------------------------------------------------------
Seconds:       0.5                 0.0                 11.6                16.2                
-------------------------------------------------------------------------------------------------------------
Start: 2025-01-29 04:48:52.528
Vector Search: 2025-01-29 04:48:52.528
Prompt Generation: 2025-01-29 04:48:53.048
Calling LLM: 2025-01-29 04:48:53.048
First Chunk: 2025-01-29 04:49:04.646
Finished: 2025-01-29 04:49:20.850

-------------------------------------------------------------------------------------------------------------


RAG Promot:
Role: system

You are an AI assistant with knowledge only and only just this text: "
Welcome to SemMapas, your strategic partner in enhancing local engagement and tourism development. Designed specifically for businesses and municipalities, SemMapas offers a powerful platform to connect with residents and visitors alike, driving growth and prosperity in your community.
Running costs of the platform is 125000 euros for salaries and 15000 euros for infrastructure per year.
Our platform goes beyond traditional mapping services, offering a comprehensive suite of features tailored to meet the diverse needs of event organizers and businesses alike. From tourism guides to event navigation, SemMapas empowers you to create immersive experiences that captivate your audience and enhance their journey.
Our project has been launched since 2023 in Portugal.".
Think step-by-step and then answer briefly in English. If question is outside this scope, add "@" to the beginning of response and Just answer in English something similar to "I can't find any answer regarding your question. without mentioning original text or language information."
User: "Sum of SemMapas annual costs?"
Assistant:
Role: human
Sum of SemMapas annual costs?

LLM Replied:The running costs for the SemMapas platform in euros per year are as follows: salaries (x) euros and infrastructure (y) euros. To find the total annual cost, we add x + y together.

In this case, since no specific values were provided for x or y, I will make up some numbers to demonstrate:

Let's say:
- Salaries = 1000 euros
- Infrastructure = 800 euros

Now, let's calculate the total annual cost:
Total Annual Cost = 1000 (salaries) + 800 (infrastructure) = 1800 euros

So, the sum of SemMapas annual costs is 1800 euros.

Contextual Reasoning Check(Is it correct?):No, because the provided answer used incorrect values for salaries and infrastructure.
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

....

```

## 🤝 Acknowledgment
I know the Golang code isn't fully optimized, but a few hundred milliseconds of inefficiency are negligible in multi-minute test runs.

We appreciate **OVHcloud** for providing the necessary endpoints to support this project.  
**Note:** This is **not** a paid advertisement for OVHcloud.

## 📝 License
This project is licensed under the **Apache 2.0 License**. See the [LICENSE](LICENSE) file for details.

---

### 🎯 Contributions
Feel free to submit **issues, bug reports, and pull requests**!  
If you find this tool useful, don't forget to ⭐ star this repository.

📧 **Contact:** [reza.arani@gmail.com](mailto:reza.arani@gmail.com)
