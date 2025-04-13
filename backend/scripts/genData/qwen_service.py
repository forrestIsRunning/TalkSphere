from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import AutoModelForCausalLM, AutoTokenizer
import torch
import uvicorn
from typing import List
import random

app = FastAPI()


# 加载模型和分词器
model_name = "Qwen/Qwen1.5-0.5B-Chat"
tokenizer = AutoTokenizer.from_pretrained(model_name, trust_remote_code=True)
model = AutoModelForCausalLM.from_pretrained(
    model_name,
    trust_remote_code=True,
    torch_dtype=torch.float32  # 使用 float32 以节省内存
)

class GenerateRequest(BaseModel):
    prompt: str
    max_new_tokens: int = 50

@app.post("/generate")
async def generate_text(request: GenerateRequest):
    try:
        # 构建完整的提示词
        messages = [
            {"role": "system", "content": "你是一个帮助生成用户简介的助手。请严格按照要求生成简介，不要添加任何引号或标点符号。"},
            {"role": "user", "content": request.prompt}
        ]
        
        # 使用模型生成文本
        text = tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=True
        )
        
        inputs = tokenizer(text, return_tensors="pt")
        outputs = model.generate(
            inputs.input_ids,
            max_new_tokens=request.max_new_tokens,
            temperature=0.9,  # 增加温度以增加随机性
            top_p=0.95,  # 增加 top_p 以允许更多样的采样
            top_k=50,  # 增加 top_k 以允许更多样的采样
            repetition_penalty=1.0,  # 降低重复惩罚以允许更多重复
            do_sample=True,  # 确保启用采样
            pad_token_id=tokenizer.pad_token_id,
            eos_token_id=tokenizer.eos_token_id
        )
        
        # 解码生成的文本，并移除提示词部分
        response = tokenizer.decode(outputs[0], skip_special_tokens=True)
        generated_text = response.split("assistant\n")[-1].strip()
        
        # 移除所有引号
        generated_text = generated_text.replace('"', '').replace('"', '')
        
        return {"generated_text": generated_text}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000) 